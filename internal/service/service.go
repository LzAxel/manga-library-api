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

	GetByID(ctx context.Context, id string) (domain.Manga, error)
	GetBySlug(ctx context.Context, slug string) (domain.Manga, error)
	Delete(ctx context.Context, userId string, mangaId string) error
	Update(ctx context.Context, userId string, roles domain.Roles, mangaDTO domain.UpdateMangaDTO) error
	UploadChapter(ctx context.Context, chapterDTO domain.UploadChapterDTO, roles domain.Roles) error
	DeleteChapter(ctx context.Context, chapterDTO domain.DeleteChapterDTO, roles domain.Roles) error
}

type Preview interface {
	Create(ctx context.Context, file multipart.File, filename string, uploaderId string) (string, error)
	Delete(ctx context.Context, previewId string) error
}

type User interface {
	GetByID(ctx context.Context, userID string) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
}

type Authorization interface {
	SignUp(ctx context.Context, userDTO domain.CreateUserDTO) error
	SignIn(ctx context.Context, userDTO domain.LoginUserDTO) (string, error)
}

type Service struct {
	Manga
	Preview
	Authorization
	User

	logger     logger.Logger
	storages   *storage.Storage
	JWTMangaer *jwt.JWTManager
	adminUser  domain.AdminUser
}

func NewService(storage *storage.Storage, JWTManager *jwt.JWTManager,
	logger logger.Logger, adminUser domain.AdminUser) *Service {

	return &Service{
		storages:      storage,
		JWTMangaer:    JWTManager,
		logger:        logger,
		Manga:         NewMangaService(storage.Manga, logger),
		Preview:       NewPreviewService(storage.Preview, logger),
		Authorization: NewAuthorizationService(storage.Authorization, logger, JWTManager, adminUser),
		User:          NewUserService(storage.User, logger),
	}
}
