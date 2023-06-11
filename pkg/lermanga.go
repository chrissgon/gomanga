package pkg

import "github.com/chrissgon/gomanga/internal/lermanga"

type lerManga manga

func (l lerManga) Search() ([]string, error) {
	return lermanga.Search(l.manga, l.chapter)
}

func NewLerManga(manga, chapter string) Provider {
	return &lerManga{
		manga:   manga,
		chapter: chapter,
	}
}
