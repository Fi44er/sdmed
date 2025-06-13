package setup_mocks

import (
	"context"
	"time"

	"github.com/Fi44er/sdmed/internal/module/auth/entity"
	"github.com/Fi44er/sdmed/internal/module/auth/pkg/constant"
	mocks "github.com/Fi44er/sdmed/internal/module/auth/usecase/auth/mock"
	"github.com/golang/mock/gomock"
)

type MockRefreshTokenDeps struct {
	User         *mocks.MockIUserUsecase
	Session      *mocks.MockISessionRepository
	TokenService *mocks.MockITokenService
}

var RefreshTokenTests = []struct {
	Name        string
	ExpectedErr error
	SetupMocks  func(ctx context.Context, m *MockRefreshTokenDeps)
}{
	{
		Name: "Success",
		SetupMocks: func(ctx context.Context, m *MockRefreshTokenDeps) {
			userFromDB := &entity.User{
				ID:       "user-id",
				Email:    "user@example.com",
				Password: "hashed123",
			}

			accessToken := "access-token"
			expiresAt := time.Now().Add(time.Hour).Unix()
			gomock.InOrder(
				m.Session.EXPECT().
					GetSessionInfo(ctx).
					Return(&entity.UserSession{UserID: "test"}, nil),

				m.TokenService.EXPECT().
					ValidateToken(gomock.Any(), gomock.Any()).
					Return(&entity.TokenDetails{UserID: "test"}, nil),

				m.User.EXPECT().
					GetByID(ctx, "test").
					Return(userFromDB, nil),

				m.TokenService.EXPECT().
					CreateToken(userFromDB.ID, gomock.Any(), gomock.Any()).
					Return(&entity.TokenDetails{
						Token:     &accessToken,
						TokenUUID: "uuid-1",
						UserID:    userFromDB.ID,
						ExpiresIn: &expiresAt,
					}, nil),
			)
		},
	},
	{
		Name:        "SessionInfoNotFound",
		ExpectedErr: constant.ErrSessionInfoNotFound,
		SetupMocks: func(ctx context.Context, m *MockRefreshTokenDeps) {
			m.Session.EXPECT().
				GetSessionInfo(ctx).
				Return(nil, constant.ErrSessionInfoNotFound)
		},
	},
	{
		Name:        "ValidateTokenError",
		ExpectedErr: constant.ErrForbidden,
		SetupMocks: func(ctx context.Context, m *MockRefreshTokenDeps) {
			m.Session.EXPECT().
				GetSessionInfo(ctx).
				Return(&entity.UserSession{UserID: "test"}, nil)

			m.TokenService.EXPECT().
				ValidateToken(gomock.Any(), gomock.Any()).
				Return(nil, constant.ErrForbidden)
		},
	},
	{
		Name:        "UserNotFound",
		ExpectedErr: constant.ErrUserNotFound,
		SetupMocks: func(ctx context.Context, m *MockRefreshTokenDeps) {
			m.Session.EXPECT().
				GetSessionInfo(ctx).
				Return(&entity.UserSession{UserID: "test"}, nil)

			m.TokenService.EXPECT().
				ValidateToken(gomock.Any(), gomock.Any()).
				Return(&entity.TokenDetails{UserID: "test"}, nil)

			m.User.EXPECT().
				GetByID(ctx, "test").
				Return(nil, constant.ErrUserNotFound)

		},
	},
	{
		Name:        "AccessTokenCreationFailed",
		ExpectedErr: constant.ErrUnprocessableEntity,
		SetupMocks: func(ctx context.Context, m *MockRefreshTokenDeps) {
			userFromDB := &entity.User{
				ID:       "user-id",
				Email:    "user@example.com",
				Password: "hashed123",
			}

			m.Session.EXPECT().
				GetSessionInfo(ctx).
				Return(&entity.UserSession{UserID: "test"}, nil)

			m.TokenService.EXPECT().
				ValidateToken(gomock.Any(), gomock.Any()).
				Return(&entity.TokenDetails{UserID: "test"}, nil)

			m.User.EXPECT().
				GetByID(ctx, "test").
				Return(userFromDB, nil)

			m.TokenService.EXPECT().
				CreateToken(userFromDB.ID, gomock.Any(), gomock.Any()).
				Return(nil, constant.ErrUnprocessableEntity)
		},
	},
}
