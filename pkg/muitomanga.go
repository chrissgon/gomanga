package pkg

import "github.com/chrissgon/gomanga/internal/muitomanga"

type muitoManga manga

func (m muitoManga) Search() ([]string, error) {
	return muitomanga.Search(m.manga, m.chapter)
}

func NewMuitoManga(manga, chapter string) Provider {
	return &muitoManga{
		manga:   manga,
		chapter: chapter,
	}
}
