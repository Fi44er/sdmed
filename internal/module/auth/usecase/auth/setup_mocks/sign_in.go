package setup_mocks

import (
	"context"
	"errors"
	"time"

	"github.com/Fi44er/sdmed/internal/module/auth/entity"
	"github.com/Fi44er/sdmed/internal/module/auth/pkg/constant"
	mocks "github.com/Fi44er/sdmed/internal/module/auth/usecase/auth/mock"

	"github.com/golang/mock/gomock"
)

type MockSignInDeps struct {
	User         *mocks.MockIUserUsecase
	Session      *mocks.MockISessionRepository
	TokenService *mocks.MockITokenService
}

var SignInTests = []struct {
	Name        string
	Input       *entity.User
	ExpectedErr error
	SetupMocks  func(ctx context.Context, m *MockSignInDeps)
}{
	{
		Name: "Success",
		Input: &entity.User{
			Email:    "user@example.com",
			Password: "pass123",
		},
		SetupMocks: func(ctx context.Context, m *MockSignInDeps) {
			userFromDB := &entity.User{
				ID:       "user-id",
				Email:    "user@example.com",
				Password: "hashed123",
			}

			accessToken := "access-token"
			refreshToken := "refresh-token"
			expiresAt := time.Now().Add(time.Hour).Unix()

			gomock.InOrder(
				m.User.EXPECT().
					GetByEmail(ctx, "user@example.com").
					Return(userFromDB, nil),

				m.User.EXPECT().
					ComparePassword(userFromDB, "pass123").
					Return(true),

				m.TokenService.EXPECT().
					CreateToken(userFromDB.ID, gomock.Any(), gomock.Any()).
					Return(&entity.TokenDetails{
						Token:     &accessToken,
						TokenUUID: "uuid-1",
						UserID:    userFromDB.ID,
						ExpiresIn: &expiresAt,
					}, nil),

				m.TokenService.EXPECT().
					CreateToken(userFromDB.ID, gomock.Any(), gomock.Any()).
					Return(&entity.TokenDetails{
						Token:     &refreshToken,
						TokenUUID: "uuid-2",
						UserID:    userFromDB.ID,
						ExpiresIn: &expiresAt,
					}, nil),

				m.Session.EXPECT().
					PutSessionInfo(ctx, gomock.Any()).
					Return(nil),
			)
		},
	},
	{
		Name: "InvalidPassword",
		Input: &entity.User{
			Email:    "user@example.com",
			Password: "wrongpass",
		},
		ExpectedErr: constant.ErrInvalidEmailOrPassword,
		SetupMocks: func(ctx context.Context, m *MockSignInDeps) {
			user := &entity.User{
				Email:    "user@example.com",
				Password: "hashed",
			}

			m.User.EXPECT().
				GetByEmail(ctx, "user@example.com").
				Return(user, nil)

			m.User.EXPECT().
				ComparePassword(user, "wrongpass").
				Return(false)
		},
	},
	{
		Name: "UserNotFound",
		Input: &entity.User{
			Email:    "notfound@example.com",
			Password: "somepass",
		},
		ExpectedErr: errors.New("user not found"),
		SetupMocks: func(ctx context.Context, m *MockSignInDeps) {
			m.User.EXPECT().
				GetByEmail(ctx, "notfound@example.com").
				Return(nil, errors.New("user not found"))
		},
	},
	{
		Name: "AccessTokenCreationFailed",
		Input: &entity.User{
			Email:    "user@example.com",
			Password: "pass123",
		},
		ExpectedErr: constant.ErrUnprocessableEntity,
		SetupMocks: func(ctx context.Context, m *MockSignInDeps) {
			user := &entity.User{
				ID:       "user-id",
				Email:    "user@example.com",
				Password: "hashed",
			}

			m.User.EXPECT().
				GetByEmail(ctx, "user@example.com").
				Return(user, nil)

			m.User.EXPECT().
				ComparePassword(user, "pass123").
				Return(true)

			m.TokenService.EXPECT().
				CreateToken(user.ID, gomock.Any(), gomock.Any()).
				Return(nil, errors.New("token error"))
		},
	},
	{
		Name: "RefreshTokenCreationFailed",
		Input: &entity.User{
			Email:    "user@example.com",
			Password: "pass123",
		},
		ExpectedErr: constant.ErrUnprocessableEntity,
		SetupMocks: func(ctx context.Context, m *MockSignInDeps) {
			user := &entity.User{
				ID:       "user-id",
				Email:    "user@example.com",
				Password: "hashed",
			}
			token := "access-token"
			expires := time.Now().Add(time.Hour).Unix()

			m.User.EXPECT().
				GetByEmail(ctx, "user@example.com").
				Return(user, nil)

			m.User.EXPECT().
				ComparePassword(user, "pass123").
				Return(true)

			m.TokenService.EXPECT().
				CreateToken(user.ID, gomock.Any(), gomock.Any()).
				Return(&entity.TokenDetails{
					Token:     &token,
					TokenUUID: "uuid-access",
					UserID:    user.ID,
					ExpiresIn: &expires,
				}, nil)

			m.TokenService.EXPECT().
				CreateToken(user.ID, gomock.Any(), gomock.Any()).
				Return(nil, errors.New("token error"))
		},
	},
	{
		Name: "SessionSaveFailed",
		Input: &entity.User{
			Email:    "user@example.com",
			Password: "pass123",
		},
		ExpectedErr: errors.New("session error"),
		SetupMocks: func(ctx context.Context, m *MockSignInDeps) {
			user := &entity.User{
				ID:       "user-id",
				Email:    "user@example.com",
				Password: "hashed",
			}
			token := "token"
			expires := time.Now().Add(time.Hour).Unix()

			m.User.EXPECT().
				GetByEmail(ctx, "user@example.com").
				Return(user, nil)

			m.User.EXPECT().
				ComparePassword(user, "pass123").
				Return(true)

			m.TokenService.EXPECT().
				CreateToken(user.ID, gomock.Any(), gomock.Any()).
				Return(&entity.TokenDetails{
					Token:     &token,
					TokenUUID: "uuid-access",
					UserID:    user.ID,
					ExpiresIn: &expires,
				}, nil)

			m.TokenService.EXPECT().
				CreateToken(user.ID, gomock.Any(), gomock.Any()).
				Return(&entity.TokenDetails{
					Token:     &token,
					TokenUUID: "uuid-refresh",
					UserID:    user.ID,
					ExpiresIn: &expires,
				}, nil)

			m.Session.EXPECT().
				PutSessionInfo(ctx, gomock.Any()).
				Return(errors.New("session error"))
		},
	},
}
