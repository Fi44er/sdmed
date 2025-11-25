package setup_mocks

import (
	"context"
	"errors"
	"time"

	auth_entity "github.com/Fi44er/sdmed/internal/module/auth/entity"
	auth_constant "github.com/Fi44er/sdmed/internal/module/auth/pkg/constant"
	mocks "github.com/Fi44er/sdmed/internal/module/auth/usecase/auth/mock"
	"github.com/golang/mock/gomock"
)

type MockForgotPasswordDeps struct {
	User   *mocks.MockIUserUsecase
	Cache  *mocks.MockICache
	Notify *mocks.MockINotificationService
}

var ForgotPasswordTests = []struct {
	Name        string
	Input       *auth_entity.Code
	ExpectedErr error
	SetupMocks  func(ctx context.Context, m *MockForgotPasswordDeps)
}{
	{
		Name: "Success",
		Input: &auth_entity.Code{
			Code:  "123456",
			Email: "user@example.com",
		},
		ExpectedErr: nil,
		SetupMocks: func(ctx context.Context, m *MockForgotPasswordDeps) {
			m.User.EXPECT().
				GetByEmail(ctx, "user@example.com").
				Return(&auth_entity.User{Email: "user@example.com"}, nil)

			m.Cache.EXPECT().
				Set(ctx, gomock.Any(), gomock.Any(), 10*time.Minute).
				Return(nil)

			m.Notify.EXPECT().
				Send(gomock.Any(), "smtp")

		},
	},
	{
		Name: "UserNotFound",
		Input: &auth_entity.Code{
			Code:  "123456",
			Email: "user@example.com",
		},
		ExpectedErr: auth_constant.ErrUserNotFound,
		SetupMocks: func(ctx context.Context, m *MockForgotPasswordDeps) {
			m.User.EXPECT().
				GetByEmail(ctx, "user@example.com").
				Return(nil, auth_constant.ErrUserNotFound)
		},
	},
	{
		Name: "CacheSetError",
		Input: &auth_entity.Code{
			Code:  "123456",
			Email: "user@example.com",
		},
		ExpectedErr: errors.New("cache set failed"),
		SetupMocks: func(ctx context.Context, m *MockForgotPasswordDeps) {
			m.User.EXPECT().
				GetByEmail(ctx, "user@example.com").
				Return(&auth_entity.User{Email: "user@example.com"}, nil)

			m.Cache.EXPECT().
				Set(ctx, gomock.Any(), gomock.Any(), 10*time.Minute).
				Return(errors.New("cache set failed"))
		},
	},
}
