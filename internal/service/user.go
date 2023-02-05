package service

import (
	"context"
	"manga-library/internal/domain"
	"manga-library/internal/storage"
	"manga-library/pkg/logger"
)

type UserService struct {
	storage storage.User
	logger  logger.Logger
}

func NewUserService(storage storage.User, logger logger.Logger) *UserService {
	return &UserService{storage: storage, logger: logger}
}

func (s *UserService) GetByID(ctx context.Context, userID string) (domain.User, error) {
	return s.storage.GetByID(ctx, userID)
}

func (s *UserService) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	return s.storage.GetByUsername(ctx, username)
}
