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

func NewMangaMemory(logger logger.Logger) *MangaMemory {
	return &MangaMemory{logger: logger}
}

func (m *MangaMemory) Create(ctx context.Context, manga domain.Manga) (string, error) {
	var storedManga []domain.Manga

	m.m.Range(func(key, value any) bool {
		mangaElem := value.(domain.Manga)
		storedManga = append(storedManga, mangaElem)

		return true
	})

	for _, m := range storedManga {
		if m.Slug == manga.Slug {
			return "", domain.ErrMangaTitleExists
		}
	}

	m.m.Store(manga.ID, manga)

	return manga.ID, nil
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

	loadedManga, ok := m.m.Load(mangaId)
	manga = loadedManga.(domain.Manga)
	if !ok {
		return manga, domain.ErrNotFound
	}

	return manga, nil
}
func (m *MangaMemory) GetBySlug(ctx context.Context, mangaSlug string) (domain.Manga, error) {
	var manga domain.Manga

	m.m.Range(func(key, value any) bool {
		storedManga := value.(domain.Manga)
		if storedManga.Slug == mangaSlug {
			manga = storedManga
			return false
		}

		return true
	})
	if manga.Title == "" {
		return manga, domain.ErrNotFound
	}

	return manga, nil
}
func (m *MangaMemory) Delete(ctx context.Context, mangaId string) error {
	m.m.Delete(mangaId)

	return nil
}
func (m *MangaMemory) Update(ctx context.Context, mangaDTO domain.UpdateMangaDTO) error {
	loadedManga, ok := m.m.Load(mangaDTO.ID)
	if !ok {
		return domain.ErrNotFound
	}
	manga := loadedManga.(domain.Manga)

	if mangaDTO.Title != nil {
		manga.Title = *mangaDTO.Title
	}
	if mangaDTO.AlternativeTitles != nil {
		manga.AlternativeTitles = *mangaDTO.AlternativeTitles
	}
	if mangaDTO.Description != nil {
		manga.Description = *mangaDTO.Description
	}
	if mangaDTO.Tags != nil {
		manga.Tags = *mangaDTO.Tags
	}
	if mangaDTO.PreviewURL != nil {
		manga.PreviewURL = *mangaDTO.PreviewURL
	}
	if mangaDTO.AgeRating != nil {
		manga.AgeRating = *mangaDTO.AgeRating
	}
	if mangaDTO.ReleaseYear != nil {
		manga.ReleaseYear = *mangaDTO.ReleaseYear
	}

	m.m.Store(manga.ID, manga)

	return nil
}

func (m *MangaMemory) UploadChapter(ctx context.Context, chapter domain.Chapter) error {
	return nil
}

func (m *MangaMemory) DeleteChapter(ctx context.Context, chapter domain.DeleteChapterDTO) error {
	return nil
}
