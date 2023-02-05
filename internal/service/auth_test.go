package service_test

import (
	"context"
	"testing"
	"time"

	"manga-library/internal/domain"
	"manga-library/internal/service"
	"manga-library/internal/storage"
	"manga-library/pkg/jwt"
	"manga-library/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func TestAuth_SignUp(t *testing.T) {
	t.Run("creating new user", func(t *testing.T) {
		var (
			l      = logger.NewLogrusLogger("info", false, false)
			jwt    = jwt.NewJWTManager("secret", int(time.Hour)*12)
			admUsr = domain.AdminUser{Username: "admin", Password: "admin"}
			svc    = service.NewAuthorizationService(storage.NewInMemoryStorage(l),
				l, jwt, admUsr)
		)

		type testCase struct {
			username    string
			password    string
			expectedErr error
		}

		testCases := []testCase{
			{
				username:    "admin",
				password:    "admin",
				expectedErr: domain.ErrUsernameExists,
			},
			{
				username:    "testuser",
				password:    "testpassword",
				expectedErr: nil,
			},
		}

		for _, tc := range testCases {
			user := domain.CreateUserDTO{
				Username: tc.username,
				Password: tc.password}

			actual := svc.SignUp(context.Background(), user)
			assert.Equal(t, tc.expectedErr, actual)
		}
	})
}

func TestAuth_SignIn(t *testing.T) {
	var (
		l       = logger.NewLogrusLogger("info", false, false)
		storage = storage.NewInMemoryStorage(l)
		jwt     = jwt.NewJWTManager("secret", int(time.Hour)*12)
		admUsr  = domain.AdminUser{Username: "admin", Password: "admin"}
		newUsr  = domain.CreateUserDTO{
			Username: "testuser",
			Password: "testpassword"}
		existedUsr = domain.CreateUserDTO{
			Username: "newuser",
			Password: "testpassword"}
		svc = service.NewAuthorizationService(storage, l, jwt, admUsr)
	)

	t.Run("login existed user", func(t *testing.T) {
		svc.SignUp(context.TODO(), existedUsr)
		token, err := svc.SignIn(context.Background(), domain.LoginUserDTO(existedUsr))
		assert.Equal(t, nil, err)
		assert.NotEmpty(t, token)
	})

	t.Run("login existed user with wrong password", func(t *testing.T) {
		user := domain.LoginUserDTO(existedUsr)
		user.Password = "2222222222222222"

		token, err := svc.SignIn(context.Background(), user)
		assert.Equal(t, domain.ErrWrongAuthCreditionals, err)
		assert.Empty(t, token)
	})

	t.Run("login unexisted user", func(t *testing.T) {
		token, err := svc.SignIn(context.Background(), domain.LoginUserDTO(newUsr))
		assert.Equal(t, domain.ErrWrongAuthCreditionals, err)
		assert.Empty(t, token)
	})

	t.Run("login admin user", func(t *testing.T) {
		token, err := svc.SignIn(context.Background(), domain.LoginUserDTO(admUsr))
		assert.Equal(t, nil, err)
		assert.NotEmpty(t, token)
	})
}
