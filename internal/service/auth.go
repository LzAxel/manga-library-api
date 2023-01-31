package service

import (
	"context"
	"manga-library/internal/domain"
	"manga-library/internal/storage"
	"manga-library/pkg/errors"
	"manga-library/pkg/hash"
	"manga-library/pkg/jwt"
	"manga-library/pkg/logger"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	hashSalt = "sajcfluyq89y414ynr9c34safwceq41c4431235c124"
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

	hashedPassword, err := hash.HashPassword(userDTO.Password, hashSalt)
	if err != nil {
		return err
	}

	id := uuid.New().String()

	user = domain.User{
		Id:           id,
		IsEditor:     false,
		Username:     userDTO.Username,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	s.logger.WithFields(logrus.Fields{"id": id, "username": userDTO.Username}).Debugln("creating user")

	return s.storage.SignUp(ctx, user)
}

func (s *AuthorizationService) SignIn(ctx context.Context, userDTO domain.LoginUserDTO) (string, error) {

	password, userId, err := s.storage.SignIn(ctx, userDTO.Username)
	if err != nil {
		return "", err
	}

	if err := hash.ComparePassword(password, userDTO.Password+hashSalt); err != nil {
		return "", errors.ErrWrongPassword
	}

	token := s.jwtManager.NewJWT(userId)

	return token, nil
}
