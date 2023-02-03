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

const (
	hashSalt = "sajcfluyq89y414ynr9c34safwceq41c4431235c124"
)

type AuthorizationService struct {
	storage    storage.Authorization
	logger     logger.Logger
	jwtManager *jwt.JWTManager
	adminUser  domain.AdminUser
}

func NewAuthorizationService(storage storage.Authorization, logger logger.Logger, jwtManager *jwt.JWTManager, adminUser domain.AdminUser) *AuthorizationService {
	return &AuthorizationService{
		storage:    storage,
		logger:     logger,
		jwtManager: jwtManager,
		adminUser:  adminUser,
	}
}

func (s *AuthorizationService) SignUp(ctx context.Context, userDTO domain.CreateUserDTO) error {
	var user domain.User

	hashedPassword, err := hash.HashPassword(userDTO.Password, hashSalt)
	if err != nil {
		return err
	}

	id := uuid.NewString()

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

	if userDTO == domain.LoginUserDTO(s.adminUser) {
		// TODO fix admin id
		token := s.jwtManager.NewJWT("admin-user-id")

		return token, nil
	}

	password, userId, err := s.storage.SignIn(ctx, userDTO.Username)
	if err != nil {
		return "", err
	}

	if err := hash.ComparePassword(password, userDTO.Password+hashSalt); err != nil {
		return "", domain.ErrWrongAuthCreditionals
	}

	token := s.jwtManager.NewJWT(userId)

	return token, nil
}
