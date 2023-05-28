package pkg

type Provider interface {
	Search() ([]string, error)
}

type NewProviderFactory func(manga, chapter string) Provider

type manga struct {
	manga   string
	chapter string
}
