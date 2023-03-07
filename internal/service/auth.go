package service

import (
	"context"
	"manga-library/internal/domain"
	"manga-library/internal/storage"
	"manga-library/pkg/hash"
	"manga-library/pkg/jwt"
	"manga-library/pkg/logger"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type AuthorizationService struct {
	storage    storage.Authorization
	logger     logger.Logger
	jwtManager *jwt.JWTManager
}

func NewAuthorizationService(storage storage.Authorization, logger logger.Logger, jwtManager *jwt.JWTManager) *AuthorizationService {
	return &AuthorizationService{
		storage:    storage,
		logger:     logger,
		jwtManager: jwtManager,
	}
}

func (s *AuthorizationService) SignUp(ctx context.Context, userDTO domain.CreateUserDTO) error {
	var user domain.User

	salt := hash.GenerateSalt()
	hashedPassword := hash.HashPassword(salt, userDTO.Password)

	id := uuid.NewString()

	user = domain.User{
		ID:           id,
		IsEditor:     false,
		IsAdmin:      false,
		Username:     userDTO.Username,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	s.logger.WithFields(logrus.Fields{"id": id, "username": userDTO.Username}).Debugln("creating user")

	return s.storage.SignUp(ctx, user)
}

func (s *AuthorizationService) SignIn(ctx context.Context, userDTO domain.LoginUserDTO) (string, error) {
	hashedPassword, userId, err := s.storage.SignIn(ctx, userDTO.Username)
	if err != nil {
		return "", domain.ErrWrongAuthCreditionals
	}

	if !hash.ComparePassword(hashedPassword, userDTO.Password) {
		return "", domain.ErrWrongAuthCreditionals
	}

	token := s.jwtManager.NewJWT(userId)

	return token, nil
}
