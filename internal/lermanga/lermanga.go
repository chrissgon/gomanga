package lermanga

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/chrissgon/gomanga/utils"
)

type mangaEntity struct {
	ID    string `json:"url"`
	Title string `json:"title"`
}

type mangasResponse struct {
	Data    []mangaEntity `json:"data"`
	Success bool          `json:"success"`
}

func Search(mangaName, strChapter string) ([]string, error) {
	mangas, err := getMangas(getMangasRequest(mangaName))

	if err != nil {
		return nil, internalError(err)
	}

	title := utils.GetTitleWithGreatestSimilarity(mangaName, getTitlesOfMangas(mangas))
	mangaID := mangas[title]

	numberPages, err := getNumberPagesOfChapter(getChapterPageRequest(mangaID, strChapter))

	if err != nil {
		return nil, internalError(err)
	}

	return getImagesURL(mangaID, strChapter, numberPages), nil
}

func getMangas(res *http.Response, err error) (map[string]string, error) {
	if err != nil {
		return nil, utils.NewError("getMangas", err)
	}

	defer res.Body.Close()

	response := mangasResponse{}
	json.NewDecoder(res.Body).Decode(&response)

	if !response.Success {
		return nil, utils.NewError("getMangas", utils.ERROR_NOT_FOUND)
	}

	response.Data = response.Data[1:]

	mangas := map[string]string{}
	for _, manga := range response.Data {
		regexMatchMangaID := regexp.MustCompile(`(https://lermanga.org/mangas/)(.*)(/)`)

		mangas[manga.Title] = regexMatchMangaID.FindStringSubmatch(manga.ID)[2]
	}

	return mangas, nil
}

func getTitlesOfMangas(mangas map[string]string) (titles []string) {
	for title := range mangas {
		titles = append(titles, title)
	}
	return
}

func getNumberPagesOfChapter(res *http.Response, err error) (int, error) {
	if err != nil {
		return 0, err
	}

	if res.StatusCode != 200 {
		return 0, utils.ERROR_NOT_FOUND
	}

	if res.Request.Response != nil && res.Request.Response.StatusCode != 200 {
		return 0, utils.ERROR_NOT_FOUND
	}

	defer res.Body.Close()

	html, err := utils.ResponseBodyToHTML(res.Body)

	if err != nil {
		return 0, err
	}

	regexMatchSelectContainer := regexp.MustCompile(`(<select class="select_paged">)(.*</option>)`)
	regexMatchSelectOptions := regexp.MustCompile(`(<option [^>]*>)`)

	strOptions := regexMatchSelectContainer.FindStringSubmatch(html)[2]
	arrOptions := regexMatchSelectOptions.FindAllString(strOptions, -1)

	return len(arrOptions) - 1, nil
}

func getImagesURL(mangaID, strChapter string, numberPages int) (images []string) {
	for i := 1; i <= numberPages; i++ {
		prefix := strings.ToUpper(string(mangaID[0]))
		url := fmt.Sprintf("https://img.lermanga.org/%s/%s/capitulo-%s/%d.jpg", prefix, mangaID, strChapter, i)
		images = append(images, url)
	}
	return
}

func internalError(err error) error {
	return utils.NewError("LERMANGA", err)
}

func getMangasRequest(mangaName string) (*http.Response, error) {
	url := fmt.Sprintf("https://lermanga.org/wp-admin/admin-ajax.php?action=wp-manga-search-manga&title=%s", mangaName)
	return utils.NewRequest(url, http.MethodGet, nil, nil)
}

func getChapterPageRequest(mangaID, strChapter string) (*http.Response, error) {
	if len(strChapter) < 2 {
		strChapter = "0" + strChapter
	}
	url := fmt.Sprintf("https://lermanga.org/capitulos/%s-capitulo-%s/", mangaID, strChapter)
	return utils.NewRequest(url, http.MethodGet, nil, nil)
}
