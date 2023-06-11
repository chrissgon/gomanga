package gomanga

import (
	"testing"
)

func TestSearchByProvider(t *testing.T) {
	_, err := SearchByProvider("one punch man", "200", LERMANGA)

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
