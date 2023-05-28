package mangalivre

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"

	"github.com/chrissgon/gomanga/utils"
)

type mangaEntity struct {
	ID    int    `json:"id_serie"`
	Title string `json:"label"`
}

type mangasResponse struct {
	Mangas []mangaEntity `json:"series"`
}

type releaseEntity struct {
	ID int `json:"id_release"`
}

type chapterEntity struct {
	Number   string                   `json:"number"`
	Releases map[string]releaseEntity `json:"releases"`
}

type chaptersResponse struct {
	Chapters []chapterEntity `json:"chapters"`
}

type imageEntity struct {
	Legacy string `json:"legacy"`
}

type imagesResponse struct {
	Images []imageEntity `json:"images"`
}

const ITEMS_PER_PAGE = 30

var REQUEST_HEADERS = map[string]string{
	"Content-Type":     "application/x-www-form-urlencoded",
	"x-requested-with": "XMLHttpRequest",
}

func Search(mangaName, strChapter string) ([]string, error) {
	mangas, err := getMangas(getMangasRequest(mangaName))

	if err != nil {
		return nil, internalError(err)
	}

	title := utils.GetTitleWithGreatestSimilarity(mangaName, getTitlesOfMangas(mangas))
	mangaID := mangas[title]

	intChapter, _ := utils.StrToInt(strChapter)
	intLastChapter, _ := utils.StrToInt(getNumberLastChapter(getNumberLastChapterRequest(mangaID)))

	if isNotReleasedChapter(intLastChapter, intChapter) {
		return nil, internalError(utils.ERROR_NOT_FOUND)
	}

	possibleNumberPages := getPossibleNumberPagesOfChapter(intChapter, intLastChapter)

	releaseID, err := getReleaseID(mangaID, strChapter, possibleNumberPages)

	if err != nil {
		return nil, internalError(err)
	}

	return getImagesByReleaseID(getImagesByReleaseIDRequest(releaseID))
}

func getMangas(res *http.Response, err error) (map[string]int, error) {
	if err != nil {
		return nil, utils.NewError("getMangas", err)
	}

	defer res.Body.Close()

	response := mangasResponse{}
	json.NewDecoder(res.Body).Decode(&response)

	if len(response.Mangas) == 0 {
		return nil, utils.NewError("getMangas", utils.ERROR_NOT_FOUND)
	}

	mangas := map[string]int{}
	for _, manga := range response.Mangas {
		mangas[manga.Title] = manga.ID
	}

	return mangas, nil
}

func getTitlesOfMangas(mangas map[string]int) (titles []string) {
	for title := range mangas {
		titles = append(titles, title)
	}
	return
}

func getNumberLastChapter(res *http.Response, err error) string {
	defer res.Body.Close()
	response := chaptersResponse{}
	json.NewDecoder(res.Body).Decode(&response)
	return response.Chapters[0].Number
}

func isNotReleasedChapter(intLastChapter int, intSearchedChapter int) bool {
	return intLastChapter < intSearchedChapter
}

func getPossibleNumberPagesOfChapter(intChapter, intLastChapter int) []int {
	numberPageOfChapter := (intChapter / ITEMS_PER_PAGE)
	numberPageOfLastChapter := (intLastChapter / ITEMS_PER_PAGE)
	numberPossiblePage := numberPageOfLastChapter - numberPageOfChapter

	if numberPossiblePage == 0 {
		return []int{1, 2, 3}
	}

	return []int{numberPossiblePage, numberPossiblePage + 1, numberPossiblePage + 2}
}

func getReleaseID(mangaID int, strChapter string, possibleNumberPages []int) (int, error) {
	timeout := make(chan error)
	go utils.TimeoutRoutine(timeout)

	chapterChan := make(chan chapterEntity)
	for _, possibleNumberPage := range possibleNumberPages {
		go func(page int) {
			res, err := getChapterByPageRequest(mangaID, page)

			chapter, err := getChapterByNumber(strChapter, getChaptersResponse(res, err))

			if err == nil {
				chapterChan <- chapter
			}
		}(possibleNumberPage)
	}

	select {
	case err := <-timeout:
		return 0, utils.NewError("getReleaseID", err)
	case response := <-chapterChan:
		for _, release := range response.Releases {
			return release.ID, nil
		}
	}

	return 0, utils.NewError("getReleaseID", utils.ERROR_NOT_FOUND)
}

func getChaptersResponse(res *http.Response, err error) chaptersResponse {
	if err != nil {
		return chaptersResponse{}
	}

	defer res.Body.Close()

	response := chaptersResponse{}
	json.NewDecoder(res.Body).Decode(&response)

	return response
}

func getChapterByNumber(strChapter string, chapters chaptersResponse) (chapterEntity, error) {
	for _, chapter := range chapters.Chapters {
		if chapter.Number == strChapter {
			return chapter, nil
		}
	}
	return chapterEntity{}, utils.NewError("getChapterByNumber", utils.ERROR_NOT_FOUND)
}

func getImagesByReleaseID(res *http.Response, err error) ([]string, error) {
	if err != nil {
		return nil, utils.NewError("getImagesByReleaseID", err)
	}

	defer res.Body.Close()

	response := imagesResponse{}
	json.NewDecoder(res.Body).Decode(&response)

	images := []string{}
	for _, image := range response.Images {
		isImageURL := regexp.MustCompile(".(gif|jpe?g|tiff?|png|webp|bmp)").MatchString(image.Legacy)

		if isImageURL {
			images = append(images, image.Legacy)
		}
	}

	return images, nil
}

func internalError(err error) error {
	return utils.NewError("MANGALIVRE", err)
}

func getMangasRequest(mangaName string) (*http.Response, error) {
	params := url.Values{}
	params.Set("search", mangaName)
	return utils.NewRequest("https://mangalivre.net/lib/search/series.json", http.MethodPost, params, REQUEST_HEADERS)
}

func getNumberLastChapterRequest(mangaID int) (*http.Response, error) {
	url := fmt.Sprintf("https://mangalivre.net/series/chapters_list.json?id_serie=%d", mangaID)
	return utils.NewRequest(url, http.MethodGet, nil, REQUEST_HEADERS)
}

func getChapterByPageRequest(mangaID, page int) (*http.Response, error) {
	url := fmt.Sprintf("https://mangalivre.net/series/chapters_list.json?id_serie=%d&page=%d", mangaID, page)
	return utils.NewRequest(url, http.MethodGet, nil, REQUEST_HEADERS)
}

func getImagesByReleaseIDRequest(releaseID int) (*http.Response, error) {
	url := fmt.Sprintf("https://mangalivre.net/leitor/pages/%v.json", releaseID)
	return utils.NewRequest(url, http.MethodGet, nil, REQUEST_HEADERS)
}
