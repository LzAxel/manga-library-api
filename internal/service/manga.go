package service

import (
	"context"
	"errors"
	"manga-library/internal/domain"
	"manga-library/internal/storage"
	"manga-library/pkg/image"
	"manga-library/pkg/logger"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gosimple/slug"

	"github.com/google/uuid"
)

// TODO: export baseUrl to env variable
const (
	uploadMangaPath = "./files/manga/"
	baseUrl         = "http://localhost:80/"
	uploadUrl       = "files/manga/"
)

type MangaService struct {
	storage storage.Manga
	logger  logger.Logger
}

func NewMangaService(storage storage.Manga, logger logger.Logger) *MangaService {
	return &MangaService{storage: storage, logger: logger}
}

func (s *MangaService) Create(ctx context.Context, userId string, mangaDTO domain.CreateMangaDTO) (string, error) {
	slug := slug.Make(mangaDTO.Title)

	for idx, tag := range mangaDTO.Tags {
		mangaDTO.Tags[idx] = strings.ToLower(tag)
	}

	if err := image.ValidateExtansion(mangaDTO.Preview.Filename); err != nil {
		return "", err
	}
	previewFile, err := mangaDTO.Preview.Open()
	if err != nil {
		return "", errors.New("failed to open file")
	}
	defer previewFile.Close()

	uploadPath := uploadMangaPath + slug

	if _, err := os.Stat(uploadPath); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(uploadPath, 0777)
		if err != nil {
			s.logger.Fatalln(err)
			return "", errors.New("failed to create directory")
		}
	}

	if err := image.UploadImage(previewFile, mangaDTO.Preview.Filename, uploadPath); err != nil {
		return "", err
	}
	manga := domain.Manga{
		ID:                uuid.NewString(),
		Title:             mangaDTO.Title,
		Author:            mangaDTO.Author,
		AlternativeTitles: mangaDTO.AlternativeTitles,
		Slug:              slug,
		Chapters:          []domain.Chapter{},
		Description:       mangaDTO.Description,
		Tags:              mangaDTO.Tags,
		PreviewURL:        uploadUrl + slug + "/preview.jpeg",
		LikeCount:         0,
		AgeRating:         mangaDTO.AgeRating,
		UploaderId:        userId,
		ReleaseYear:       mangaDTO.ReleaseYear,
		IsPublished:       false,
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

func (s *MangaService) GetByTags(ctx context.Context, tags []string, offset int) ([]domain.Manga, error) {
	for idx, tag := range tags {
		tags[idx] = strings.ToLower(tag)
	}
	s.logger.Debugln(tags)
	return s.storage.GetByTags(ctx, tags, offset)
}

func (s *MangaService) Delete(ctx context.Context, userId string, mangaId string) error {
	return s.storage.Delete(ctx, mangaId)
}

func (s *MangaService) Update(ctx context.Context, userId string, roles domain.Roles, mangaDTO domain.UpdateMangaDTO) error {
	manga, err := s.storage.GetById(ctx, mangaDTO.ID)
	if err != nil {
		return err
	}
	if manga.UploaderId != userId && !roles.IsAdmin && !roles.IsEditor {
		return domain.ErrNotTheOwner
	}
	if mangaDTO.IsPublished != nil && !roles.IsAdmin {
		return domain.ErrMangaPublishByUser
	}
	if mangaDTO.Preview != nil {
		if err := image.ValidateExtansion(mangaDTO.Preview.Filename); err != nil {
			return err
		}
		previewFile, err := mangaDTO.Preview.Open()
		if err != nil {
			return errors.New("failed to open file")
		}
		defer previewFile.Close()

		previewPath := filepath.Join(uploadMangaPath, manga.Slug)

		if err := image.UploadImage(previewFile, "preview.jpeg", previewPath); err != nil {
			return err
		}
	}

	if mangaDTO.Title != nil {
		mangaDTO.Slug = slug.Make(*mangaDTO.Title)
	}

	if err := s.storage.Update(ctx, mangaDTO); err != nil {
		return err
	}

	return nil
}
