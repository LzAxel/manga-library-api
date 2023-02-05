package storage

import (
	"context"
	"manga-library/internal/domain"
	"manga-library/internal/storage/memory"
	"manga-library/internal/storage/mongodb"
	"manga-library/pkg/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

type Manga interface {
	Create(ctx context.Context, manga domain.Manga) (string, error)
	GetLatest(ctx context.Context) ([]domain.Manga, error)
	GetById(ctx context.Context, mangaId string) (domain.Manga, error)
	GetBySlug(ctx context.Context, mangaSlug string) (domain.Manga, error)
	Delete(ctx context.Context, mangaId string) error
	Update(ctx context.Context, mangaDTO domain.UpdateMangaDTO) error
}

type Preview interface {
	Create(ctx context.Context, preview domain.Preview) (string, error)
	Delete(ctx context.Context, previewId string) error
}

type User interface {
	GetByID(ctx context.Context, userID string) (domain.User, error)
	GetByUsername(ctx context.Context, username string) (domain.User, error)
}

type Authorization interface {
	SignUp(ctx context.Context, user domain.User) error
	SignIn(ctx context.Context, username string) (password, userId string, err error)
}

type Storage struct {
	Manga
	Preview
	Authorization
	User

	logger logger.Logger
}

func NewStorage(db *mongo.Database, logger logger.Logger) *Storage {
	return &Storage{
		logger:        logger,
		Manga:         mongodb.NewMangaMongoDB(db),
		Preview:       mongodb.NewPreviewMongoDB(db),
		Authorization: mongodb.NewAuthorizationMongoDB(logger, db),
		User:          mongodb.NewUserMongoDB(db),
	}
}

func NewInMemoryStorage(logger logger.Logger) *Storage {
	return &Storage{
		logger:        logger,
		Manga:         memory.NewMangaMemory(logger),
		Preview:       memory.NewPreviewMemory(logger),
		Authorization: memory.NewAuthMemory(logger),
	}
}
