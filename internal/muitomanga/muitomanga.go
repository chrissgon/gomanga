package muitomanga

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/chrissgon/gomanga/utils"
)

type mangaEntity struct {
	ID    string `json:"url"`
	Title string `json:"titulo"`
}

type mangasResponse []mangaEntity

var REQUEST_HEADERS = map[string]string{
	"x-requested-with": "XMLHttpRequest",
}

func Search(mangaName, strChapter string) ([]string, error) {
	mangas, err := getMangas(getMangasRequest(mangaName))

	if err != nil {
		return nil, internalError(err)
	}

	title := utils.GetTitleWithGreatestSimilarity(mangaName, getTitlesOfMangas(mangas))
	mangaID := mangas[title]

	html, err := getPageOfChapter(mangaID, strChapter)

	if err != nil {
		return nil, internalError(err)
	}

	return getImagesByHTML(html)
}

func getMangas(res *http.Response, err error) (map[string]string, error) {
	if err != nil {
		return nil, utils.NewError("getMangas", err)
	}

	defer res.Body.Close()

	response := mangasResponse{}
	json.NewDecoder(res.Body).Decode(&response)

	if len(response) == 0 {
		return nil, utils.NewError("getMangas", utils.ERROR_NOT_FOUND)
	}

	mangas := map[string]string{}
	for _, manga := range response {
		mangas[manga.Title] = manga.ID
	}

	return mangas, nil
}

func getTitlesOfMangas(mangas map[string]string) (titles []string) {
	for title := range mangas {
		titles = append(titles, title)
	}
	return
}

func getPageOfChapter(mangaID, strChapter string) (string, error) {
	url := fmt.Sprintf("https://muitomanga.com/ler/%s/capitulo-%s", mangaID, strChapter)
	res, err := utils.NewRequest(url, http.MethodGet, nil, nil)

	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		return "", utils.ERROR_NOT_FOUND
	}

	defer res.Body.Close()

	return utils.ResponseBodyToHTML(res.Body)
}

func getImagesByHTML(html string) ([]string, error) {
	regexMatchImagesURL := regexp.MustCompile(`\["https:(.*)\]`)

	strImagesURL := regexMatchImagesURL.FindString(html)

	images := []string{}
	err := json.Unmarshal([]byte(strImagesURL), &images)

	return images, err
}

func internalError(err error) error {
	return utils.NewError("MUITOMANGA", err)
}

func getMangasRequest(mangaName string) (*http.Response, error) {
	url := fmt.Sprintf("https://muitomanga.com/lib/search-ajax-json?q=%s", mangaName)
	return utils.NewRequest(url, http.MethodGet, nil, REQUEST_HEADERS)
}
