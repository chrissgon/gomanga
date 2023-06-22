package lermanga

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

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

	return getImagesURL(getChapterPageRequest(mangaID, strChapter))
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

func getImagesURL(res *http.Response, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, utils.ERROR_NOT_FOUND
	}

	if res.Request.Response != nil && res.Request.Response.StatusCode != 200 {
		return nil, utils.ERROR_NOT_FOUND
	}

	defer res.Body.Close()

	html, err := utils.ResponseBodyToHTML(res.Body)

	if err != nil {
		return nil, err
	}

	regexMatchScriptURL := regexp.MustCompile(`(https://lermanga.org/wp-content/litespeed/js/)(.*)(\.js)`)
	regexMatchJSON := regexp.MustCompile(`(?mi)(imagens_cap=)(.*"])`)

	url := regexMatchScriptURL.FindString(html)

	res, _ = http.Get(url)

	js, err := utils.ResponseBodyToHTML(res.Body)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	parts := regexMatchJSON.FindStringSubmatch(js)

	if len(parts) < 2 {
		return nil, utils.ERROR_NOT_FOUND
	}

	var images []string
	json.Unmarshal([]byte(parts[2]), &images)

	return images, nil
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
