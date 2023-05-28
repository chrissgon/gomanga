package muitomanga

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/chrissgon/gomanga/utils"
)

func TestGetMangas(t *testing.T) {
	expect := map[string]string{
		"Black Clover":                         "black-clover",
		"Black Clover Gaiden: Quartet Knights": "black-clover-gaiden-quartet-knights",
	}
	have, err := getMangas(utils.MakeMockRequest("./mocks/BLACK_CLOVER.json"))

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(utils.FormatTestError(expect, have))
	}
}

func TestGetTitlesOfMangas(t *testing.T) {
	mangas := map[string]string{
		"Black Clover":                         "black-clover",
		"Black Clover Gaiden: Quartet Knights": "black-clover-gaiden-quartet-knights",
	}

	expect := []string{"Black Clover", "Black Clover Gaiden: Quartet Knights"}
	have := getTitlesOfMangas(mangas)

	if len(expect) != len(have) {
		t.Errorf(utils.FormatTestError(expect, have))
	}
}

func TestGetImagesByHTML(t *testing.T) {
	bytes, err := ioutil.ReadFile("./mocks/BLACK_CLOVER_CHAPTER_360.html")

	if err != nil {
		t.Error(err)
	}

	expect := []string{
		"https://imgs2.muitomanga.com/imgs/black-clover/360/1.jpg",
		"https://imgs2.muitomanga.com/imgs/black-clover/360/2.jpg",
		"https://imgs2.muitomanga.com/imgs/black-clover/360/3.jpg",
		"https://imgs2.muitomanga.com/imgs/black-clover/360/4.jpg",
		"https://imgs2.muitomanga.com/imgs/black-clover/360/5.jpg",
		"https://imgs2.muitomanga.com/imgs/black-clover/360/6.jpg",
		"https://imgs2.muitomanga.com/imgs/black-clover/360/7.jpg",
		"https://imgs2.muitomanga.com/imgs/black-clover/360/8.jpg",
		"https://imgs2.muitomanga.com/imgs/black-clover/360/9.jpg",
		"https://imgs2.muitomanga.com/imgs/black-clover/360/10.jpg",
		"https://imgs2.muitomanga.com/imgs/black-clover/360/11.jpg",
		"https://imgs2.muitomanga.com/imgs/black-clover/360/12.jpg",
		"https://imgs2.muitomanga.com/imgs/black-clover/360/13.jpg",
		"https://imgs2.muitomanga.com/imgs/black-clover/360/14.jpg",
		"https://imgs2.muitomanga.com/imgs/black-clover/360/15.jpg",
	}
	have, err := getImagesByHTML(string(bytes))

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(expect, have) {
		t.Errorf(utils.FormatTestError(expect, have))
	}
}
