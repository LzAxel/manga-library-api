package memory

import (
	"context"
	"errors"
	"manga-library/internal/domain"
	"manga-library/pkg/logger"
	"sync"
	"time"

	"github.com/google/uuid"
)

type AuthMemory struct {
	m      sync.Map
	logger logger.Logger
}

func NewAuthMemory(logger logger.Logger) *AuthMemory {
	return &AuthMemory{logger: logger}
}

func (m *AuthMemory) SignUp(ctx context.Context, user domain.User) error {
	m.logger.Debugln("storage signing up accont")
	user.ID = uuid.NewString()
	user.CreatedAt = time.Now()
	if _, ok := m.m.Load(user.ID); ok {
		return errors.New("username already exists")
	}

	m.m.Store(user.Username, user)

	return nil
}

func (m *AuthMemory) SignIn(ctx context.Context, username string) (password, userId string, err error) {
	m.logger.Debugf("getting from database: %s", username)
	var user domain.User

	loadedUser, ok := m.m.Load(username)
	if !ok {
		return "", "", domain.ErrNotFound
	}
	user, ok = loadedUser.(domain.User)
	if !ok {
		return "", "", domain.ErrFailedToGet
	}

	return user.PasswordHash, user.ID, nil
}
