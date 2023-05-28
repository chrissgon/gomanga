package mangalivre

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/chrissgon/gomanga/utils"
)

func TestGetMangas(t *testing.T) {
	expect := map[string]int{"Black Clover": 1751, "Black Clover Gaiden: Quartet Knights": 9284, "Black Clover SD: Asta-kun Mahou Tei e no Michi": 12553}
	have, err := getMangas(utils.MakeMockRequest("./mocks/BLACK_CLOVER.json"))

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(utils.FormatTestError(expect, have))
	}
}

func TestGetTitlesOfMangas(t *testing.T) {
	mangas := map[string]int{"Black Clover": 1751, "Black Clover Gaiden: Quartet Knights": 9284, "Black Clover SD: Asta-kun Mahou Tei e no Michi": 12553}

	expect := []string{"Black Clover", "Black Clover Gaiden: Quartet Knights", "Black Clover SD: Asta-kun Mahou Tei e no Michi"}
	have := getTitlesOfMangas(mangas)

	if len(expect) != len(have) {
		t.Errorf(utils.FormatTestError(expect, have))
	}
}

func TestGetNumberLastChapter(t *testing.T) {
	expect := "360"
	have := getNumberLastChapter(utils.MakeMockRequest("./mocks/BLACK_CLOVER_CHAPTERS.json"))

	if expect != have {
		t.Errorf(utils.FormatTestError(expect, have))
	}
}

func TestIsNotReleasedChapter(t *testing.T) {
	expect := true
	have := isNotReleasedChapter(1, 2)

	if expect != have {
		t.Errorf(utils.FormatTestError(expect, have))
	}

	expect = false
	have = isNotReleasedChapter(2, 1)

	if expect != have {
		t.Errorf(utils.FormatTestError(expect, have))
	}
}

func TestGetPossibleNumberPagesOfChapter(t *testing.T) {
	expect := []int{5, 6, 7}
	have := getPossibleNumberPagesOfChapter(548, 700)

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(utils.FormatTestError(expect, have))
	}

	expect = []int{1, 2, 3}
	have = getPossibleNumberPagesOfChapter(700, 700)

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(utils.FormatTestError(expect, have))
	}
}

func TestGetChaptersResponse(t *testing.T) {
	res, err := utils.MakeMockRequest("./mocks/BLACK_CLOVER_CHAPTERS.json")

	if err != nil {
		t.Error(err)
	}

	expect := chaptersResponse{}
	json.NewDecoder(res.Body).Decode(&expect)

	have := getChaptersResponse(utils.MakeMockRequest("./mocks/BLACK_CLOVER_CHAPTERS.json"))

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(utils.FormatTestError(expect, have))
	}
}

func TestGetImagesByReleaseID(t *testing.T) {
	expect := []string{
		"https://static2.mangalivre.net/firefox/kgEsZfh43BotHF83nemP_g/m12015740/1751/441727/466886/0.jpg",
		"https://static2.mangalivre.net/firefox/R1fXoT8Pohc-vp8YYucu0A/m12015739/1751/441727/466886/00.jpg",
		"https://static2.mangalivre.net/firefox/Ls2hH1MhGBQcJFdQQHFuYQ/m12015741/1751/441727/466886/1.png",
		"https://static2.mangalivre.net/firefox/bJVRTDQRqOilV1l8_6MdLg/m12015742/1751/441727/466886/2.png",
		"https://static2.mangalivre.net/firefox/dwKPQqH1lWCM5Yl6-3wswA/m12015743/1751/441727/466886/3.png",
		"https://static2.mangalivre.net/firefox/CAprzMC_WDkWjs6691hL9Q/m12015744/1751/441727/466886/4.png",
		"https://static2.mangalivre.net/firefox/I1NnvMnSt-wSP8xPJb-wBw/m12015745/1751/441727/466886/5.png",
		"https://static2.mangalivre.net/firefox/dDfKjPiQopYqvtUQZ68QdA/m12015746/1751/441727/466886/6.png",
		"https://static2.mangalivre.net/firefox/Xmh2AwLtia2OzAzmeZ6tIQ/m12015747/1751/441727/466886/7.png",
		"https://static2.mangalivre.net/firefox/7wxSJPYyMPIr7g_6H2CUYQ/m12015748/1751/441727/466886/8.png",
		"https://static2.mangalivre.net/firefox/ytJd4iyg4O7o9yq82_Bq3A/m12015749/1751/441727/466886/9.png",
		"https://static2.mangalivre.net/firefox/RJTfDr3g7r65VRi293WlkQ/m12015750/1751/441727/466886/10.png",
		"https://static2.mangalivre.net/firefox/_3u0ESrbsbJwl2ojeMQhMA/m12015751/1751/441727/466886/11.png",
		"https://static2.mangalivre.net/firefox/HdHekyfnIDEFa4tXt36z7Q/m12015752/1751/441727/466886/12.png",
		"https://static2.mangalivre.net/firefox/1t1eLOO8n51pG6_XmI3vLQ/m12015753/1751/441727/466886/13.png",
		"https://static2.mangalivre.net/firefox/m1x98A0EeUNZxqc7UnszzQ/m12015754/1751/441727/466886/14.jpg",
	}

	have, err := getImagesByReleaseID(utils.MakeMockRequest("./mocks/BLACK_CLOVER_PAGES.json"))

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(utils.FormatTestError(expect, have))
	}
}
