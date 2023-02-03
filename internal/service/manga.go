package service

import (
	"context"
	"errors"
	"manga-library/internal/domain"
	"manga-library/internal/storage"
	"manga-library/pkg/logger"
	"time"

	"github.com/gosimple/slug"

	"github.com/google/uuid"
)

type MangaService struct {
	storage storage.Manga
	logger  logger.Logger
}

func NewMangaService(storage storage.Manga, logger logger.Logger) *MangaService {
	return &MangaService{storage: storage, logger: logger}
}

func (s *MangaService) Create(ctx context.Context, userId string, mangaDTO domain.CreateMangaDTO) (string, error) {
	manga := domain.Manga{
		Id:                uuid.NewString(),
		Title:             mangaDTO.Title,
		AlternativeTitles: mangaDTO.AlternativeTitles,
		Slug:              slug.Make(mangaDTO.Title),
		Description:       mangaDTO.Description,
		Tags:              mangaDTO.Tags,
		PreviewURL:        mangaDTO.PreviewURL,
		LikeCount:         0,
		AgeRating:         mangaDTO.AgeRating,
		UploaderId:        userId,
		ReleaseYear:       mangaDTO.ReleaseYear,
		CreatedAt:         time.Now(),
	}

	return s.storage.Create(ctx, manga)
}

func (s *MangaService) GetLatest(ctx context.Context) ([]domain.Manga, error) {
	return s.storage.GetLatest(ctx)
}

func (s *MangaService) GetByID(ctx context.Context, id string) (domain.Manga, error) {
	return s.storage.GetById(ctx, id)
}

func (s *MangaService) GetBySlug(ctx context.Context, slug string) (domain.Manga, error) {
	return s.storage.GetBySlug(ctx, slug)
}

func (s *MangaService) Delete(ctx context.Context, userId string, mangaId string) error {
	return s.storage.Delete(ctx, mangaId)
}

func (s *MangaService) Update(ctx context.Context, userId string, mangaDTO domain.UpdateMangaDTO) error {
	manga, err := s.storage.GetById(ctx, mangaDTO.Id)
	if err != nil {
		return err
	}

	if manga.UploaderId != userId {
		// TODO: make general error in utils for this case
		return errors.New("you are not an owner")
	}

	if mangaDTO.Title != nil {
		mangaDTO.Slug = slug.Make(*mangaDTO.Title)
	}

	if err := s.storage.Update(ctx, mangaDTO); err != nil {
		return err
	}

	return nil
}
