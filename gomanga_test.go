package gomanga

import (
	"testing"
)

func TestSearchByProvider(t *testing.T) {
	_, err := SearchByProvider("jujutsu", "3", LERMANGA)

	if err != nil {
		t.Error(err)
	}
}

func TestSearchByProviders(t *testing.T) {
	_, err := SearchByProviders("vinland", "1")

	if err != nil {
		t.Error(err)
	}
}
