package memory

import (
	"context"
	"manga-library/internal/domain"
	"manga-library/pkg/logger"
	"sort"
	"sync"
)

type MangaMemory struct {
	m      sync.Map
	logger logger.Logger
}

// TODO: make unified exists error
func NewMangaMemory(logger logger.Logger) *MangaMemory {
	return &MangaMemory{logger: logger}
}

func (m *MangaMemory) Create(ctx context.Context, manga domain.Manga) (string, error) {
	var mangaId string

	m.m.Store(manga.Id, manga)

	return mangaId, nil
}
func (m *MangaMemory) GetLatest(ctx context.Context) ([]domain.Manga, error) {
	var mangaList []domain.Manga
	var allMangaList []domain.Manga

	m.m.Range(func(key, value any) bool {
		mangaElem := value.(domain.Manga)
		allMangaList = append(allMangaList, mangaElem)

		return true
	})

	sort.Slice(allMangaList, func(i, j int) bool {
		return allMangaList[i].CreatedAt.After(allMangaList[j].CreatedAt)
	})

	if len(allMangaList) >= 20 {
		mangaList = allMangaList[:20]
	} else {
		mangaList = allMangaList
	}

	return mangaList, nil
}

func (m *MangaMemory) GetById(ctx context.Context, mangaId string) (domain.Manga, error) {
	var manga domain.Manga

	return manga, nil
}
func (m *MangaMemory) GetBySlug(ctx context.Context, mangaSlug string) (domain.Manga, error) {
	var manga domain.Manga

	return manga, nil
}
func (m *MangaMemory) Delete(ctx context.Context, mangaId string) error {
	return nil
}
func (m *MangaMemory) Update(ctx context.Context, mangaDTO domain.UpdateMangaDTO) error {
	return nil
}
