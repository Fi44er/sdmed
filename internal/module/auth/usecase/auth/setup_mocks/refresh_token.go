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
			userFromDB := &auth_entity.User{
				ID:       "user-id",
				Email:    "user@example.com",
				Password: "hashed123",
			}

			accessToken := "access-token"
			expiresAt := time.Now().Add(time.Hour).Unix()
			gomock.InOrder(
				m.Session.EXPECT().
					GetSessionInfo(ctx).
					Return(&auth_entity.UserSession{UserID: "test"}, nil),

				m.TokenService.EXPECT().
					ValidateToken(gomock.Any(), gomock.Any()).
					Return(&auth_entity.TokenDetails{UserID: "test"}, nil),

				m.User.EXPECT().
					GetByID(ctx, "test").
					Return(userFromDB, nil),

				m.TokenService.EXPECT().
					CreateToken(userFromDB.ID, gomock.Any(), gomock.Any()).
					Return(&auth_entity.TokenDetails{
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
		ExpectedErr: auth_constant.ErrSessionInfoNotFound,
		SetupMocks: func(ctx context.Context, m *MockRefreshTokenDeps) {
			m.Session.EXPECT().
				GetSessionInfo(ctx).
				Return(nil, auth_constant.ErrSessionInfoNotFound)
		},
	},
	{
		Name:        "ValidateTokenError",
		ExpectedErr: auth_constant.ErrForbidden,
		SetupMocks: func(ctx context.Context, m *MockRefreshTokenDeps) {
			m.Session.EXPECT().
				GetSessionInfo(ctx).
				Return(&auth_entity.UserSession{UserID: "test"}, nil)

			m.TokenService.EXPECT().
				ValidateToken(gomock.Any(), gomock.Any()).
				Return(nil, auth_constant.ErrForbidden)
		},
	},
	{
		Name:        "UserNotFound",
		ExpectedErr: auth_constant.ErrUserNotFound,
		SetupMocks: func(ctx context.Context, m *MockRefreshTokenDeps) {
			m.Session.EXPECT().
				GetSessionInfo(ctx).
				Return(&auth_entity.UserSession{UserID: "test"}, nil)

			m.TokenService.EXPECT().
				ValidateToken(gomock.Any(), gomock.Any()).
				Return(&auth_entity.TokenDetails{UserID: "test"}, nil)

			m.User.EXPECT().
				GetByID(ctx, "test").
				Return(nil, auth_constant.ErrUserNotFound)

		},
	},
	{
		Name:        "AccessTokenCreationFailed",
		ExpectedErr: auth_constant.ErrUnprocessableEntity,
		SetupMocks: func(ctx context.Context, m *MockRefreshTokenDeps) {
			userFromDB := &auth_entity.User{
				ID:       "user-id",
				Email:    "user@example.com",
				Password: "hashed123",
			}

			m.Session.EXPECT().
				GetSessionInfo(ctx).
				Return(&auth_entity.UserSession{UserID: "test"}, nil)

			m.TokenService.EXPECT().
				ValidateToken(gomock.Any(), gomock.Any()).
				Return(&auth_entity.TokenDetails{UserID: "test"}, nil)

			m.User.EXPECT().
				GetByID(ctx, "test").
				Return(userFromDB, nil)

			m.TokenService.EXPECT().
				CreateToken(userFromDB.ID, gomock.Any(), gomock.Any()).
				Return(nil, auth_constant.ErrUnprocessableEntity)
		},
	},
}
