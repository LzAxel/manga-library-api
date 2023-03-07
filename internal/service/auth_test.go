package service_test

import (
	"context"
	"testing"
	"time"

	"manga-library/internal/domain"
	"manga-library/internal/service"
	mock_storage "manga-library/internal/storage/mocks"
	"manga-library/pkg/hash"
	"manga-library/pkg/jwt"
	"manga-library/pkg/logger"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuth_SignUp(t *testing.T) {
	var (
		l    = logger.NewLogrusLogger("info", false, false)
		jwt  = jwt.NewJWTManager("secret", int(time.Hour)*12)
		ctrl = gomock.NewController(t)
		stor = mock_storage.NewMockAuthorization(ctrl)
		svc  = service.NewAuthorizationService(stor, l, jwt)
		ctx  = context.Background()
	)

	t.Run("sing up existed user", func(t *testing.T) {
		user := domain.CreateUserDTO{
			Username: "existedUser",
			Password: "existedPassword",
		}
		stor.EXPECT().SignUp(ctx, gomock.Any()).Return(domain.ErrUsernameExists)
		actual := svc.SignUp(ctx, user)

		assert.ErrorIs(t, actual, domain.ErrUsernameExists)
	})

	t.Run("sing up new user", func(t *testing.T) {
		user := domain.CreateUserDTO{
			Username: "newUser",
			Password: "newPassword",
		}
		stor.EXPECT().SignUp(ctx, gomock.Any()).Return(nil)
		actual := svc.SignUp(context.TODO(), user)

		require.NoError(t, actual)
	})
}

func TestAuth_SignIn(t *testing.T) {
	var (
		l            = logger.NewLogrusLogger("info", false, false)
		jwt          = jwt.NewJWTManager("secret", int(time.Hour)*12)
		ctrl         = gomock.NewController(t)
		stor         = mock_storage.NewMockAuthorization(ctrl)
		salt         = hash.GenerateSalt()
		passwordHash = hash.HashPassword(salt, "password")
		svc          = service.NewAuthorizationService(stor, l, jwt)
		ctx          = context.Background()
	)

	type testCase struct {
		name    string
		data    domain.LoginUserDTO
		prepare func(username string)

		expErr    error
		expReturn bool
	}

	testCases := []testCase{
		{
			name: "login unexisted user",
			data: domain.LoginUserDTO{
				Username: "unexistedUser",
				Password: "password",
			},
			prepare: func(username string) {
				stor.EXPECT().SignIn(ctx, gomock.Any()).Return(passwordHash, "user-id", domain.ErrNotFound)
			},
			expErr:    domain.ErrWrongAuthCreditionals,
			expReturn: false,
		},
		{
			name: "login existed user",
			data: domain.LoginUserDTO{
				Username: "existedUser",
				Password: "password",
			},
			prepare: func(username string) {
				stor.EXPECT().SignIn(ctx, username).Return(passwordHash, "user-id", nil)
			},
			expErr:    nil,
			expReturn: true,
		},
		{
			name: "login with invalid password",
			data: domain.LoginUserDTO{
				Username: "existedUser",
				Password: "invalidPassword",
			},
			prepare: func(username string) {
				stor.EXPECT().SignIn(ctx, username).Return(passwordHash, "user-id", nil)
			},
			expErr:    domain.ErrWrongAuthCreditionals,
			expReturn: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.prepare != nil {
				testCase.prepare(testCase.data.Username)
			}
			actual, err := svc.SignIn(ctx, testCase.data)
			assert.Equal(t, testCase.expErr, err)
			if testCase.expReturn {
				assert.NotEmpty(t, actual)
			}
		})
	}
}
