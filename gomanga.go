package gomanga

import (
	"github.com/chrissgon/gomanga/pkg"
	"github.com/chrissgon/gomanga/utils"
)

type PROVIDER_NAMES string

const (
	MANGALIVRE PROVIDER_NAMES = "MANGALIVRE"
	MUITOMANGA                = "MUITOMANGA"
)

var factories = map[string]pkg.NewProviderFactory{
	"MANGALIVRE": pkg.NewMangaLivre,
	"MUITOMANGA": pkg.NewMuitoManga,
}

// make provider
func NewProvider(manga, chapter string, providerName PROVIDER_NAMES) pkg.Provider {
	return factories[string(providerName)](manga, chapter)
}

// make providers
func NewProviders(manga, chapter string) (all []pkg.Provider) {
	for _, fs := range factories {
		all = append(all, fs(manga, chapter))
	}
	return
}

// search manga by especific provider
func SearchByProvider(manga, chapter string, providerName PROVIDER_NAMES) ([]string, error) {
	return NewProvider(manga, chapter, providerName).Search()
}

// search manga by providers
func SearchByProviders(manga, chapter string) ([]string, error) {
	providers := NewProviders(manga, chapter)
	calls := len(providers)

	response := make(chan interface{})
	for _, provider := range providers {
		go asyncSearch(provider, response)
	}

	for value := range response {
		calls--

		switch value.(type) {
		case []string:
			return value.([]string), nil
		default:
			if calls == 0 {
				return []string{}, utils.NewError("SearchByProviders", value.(error))
			}
		}

	}

	return []string{}, utils.NewError("SearchByProviders", utils.ERROR_NOT_FOUND)
}

func asyncSearch(provider pkg.Provider, response chan interface{}) {
	url, err := provider.Search()

	if err != nil {
		response <- err
		return
	}

	response <- url
}
