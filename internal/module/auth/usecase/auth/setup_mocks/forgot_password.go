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

type MockForgotPasswordDeps struct {
	User   *mocks.MockIUserUsecase
	Cache  *mocks.MockICache
	Notify *mocks.MockINotificationService
}

var ForgotPasswordTests = []struct {
	Name        string
	Input       *entity.Code
	ExpectedErr error
	SetupMocks  func(ctx context.Context, m *MockForgotPasswordDeps)
}{
	{
		Name: "Success",
		Input: &entity.Code{
			Code:  "123456",
			Email: "user@example.com",
		},
		ExpectedErr: nil,
		SetupMocks: func(ctx context.Context, m *MockForgotPasswordDeps) {
			m.User.EXPECT().
				GetByEmail(ctx, "user@example.com").
				Return(&entity.User{Email: "user@example.com"}, nil)

			m.Cache.EXPECT().
				Set(ctx, gomock.Any(), gomock.Any(), 10*time.Minute).
				Return(nil)

			m.Notify.EXPECT().
				Send(gomock.Any(), "smtp")

		},
	},
	{
		Name: "UserNotFound",
		Input: &entity.Code{
			Code:  "123456",
			Email: "user@example.com",
		},
		ExpectedErr: constant.ErrUserNotFound,
		SetupMocks: func(ctx context.Context, m *MockForgotPasswordDeps) {
			m.User.EXPECT().
				GetByEmail(ctx, "user@example.com").
				Return(nil, constant.ErrUserNotFound)
		},
	},
	{
		Name: "CacheSetError",
		Input: &entity.Code{
			Code:  "123456",
			Email: "user@example.com",
		},
		ExpectedErr: errors.New("cache set failed"),
		SetupMocks: func(ctx context.Context, m *MockForgotPasswordDeps) {
			m.User.EXPECT().
				GetByEmail(ctx, "user@example.com").
				Return(&entity.User{Email: "user@example.com"}, nil)

			m.Cache.EXPECT().
				Set(ctx, gomock.Any(), gomock.Any(), 10*time.Minute).
				Return(errors.New("cache set failed"))
		},
	},
}
