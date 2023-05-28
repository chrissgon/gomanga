package pkg

import "github.com/chrissgon/gomanga/internal/mangalivre"

type mangaLivre manga

func (m mangaLivre) Search() ([]string, error) {
	return mangalivre.Search(m.manga, m.chapter)
}

func NewMangaLivre(manga, chapter string) Provider {
	return &mangaLivre{
		manga:   manga,
		chapter: chapter,
	}
}
