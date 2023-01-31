package service

import (
	"context"
	"manga-library/internal/domain"
	"manga-library/internal/storage"
	"manga-library/pkg/jwt"
	"manga-library/pkg/logger"
	"mime/multipart"
)

type Manga interface {
	Create(ctx context.Context, userId string, mangaDTO domain.CreateMangaDTO) (string, error)
	GetLatest(ctx context.Context) ([]domain.Manga, error)

	Get(ctx context.Context, options domain.GetMangaDTO) (domain.Manga, error)
	Delete(ctx context.Context, userId string, mangaId string) error
	Update(ctx context.Context, userId string, mangaDTO domain.UpdateMangaDTO) error
}

type Preview interface {
	Create(ctx context.Context, file multipart.File, filename string, uploaderId string) (string, error)
	Delete(ctx context.Context, previewId string) error
}

type Authorization interface {
	SignUp(ctx context.Context, userDTO domain.CreateUserDTO) error
	SignIn(ctx context.Context, userDTO domain.LoginUserDTO) (string, error)
}

type Service struct {
	Manga
	Preview
	Authorization

	logger     logger.Logger
	storages   *storage.Storage
	JWTMangaer *jwt.JWTManager
}

func NewService(storage *storage.Storage, JWTManager *jwt.JWTManager, logger logger.Logger) *Service {
	return &Service{
		storages:      storage,
		JWTMangaer:    JWTManager,
		logger:        logger,
		Manga:         NewMangaService(storage.Manga, logger),
		Preview:       NewPreviewService(storage.Preview, logger),
		Authorization: NewAuthorizationService(storage.Authorization, logger, JWTManager),
	}
}
