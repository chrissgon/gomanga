package utils

import "testing"

func TestGetTitleWithGreatestSimilarity(t *testing.T) {
	expect := "Black Clover"
	have := GetTitleWithGreatestSimilarity("black clover", []string{"Black Clover", "Black Clover Gaiden: Quartet Knights", "Black Clover SD: Asta-kun Mahou Tei e no Michi"})

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}

func TestStrToInt(t *testing.T) {
	expect := 10
	have, _ := StrToInt("10.5")

	if expect != have {
		t.Errorf(FormatTestError(expect, have))
	}
}
