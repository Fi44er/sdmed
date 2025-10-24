package setup_mocks

import (
	"context"

	mocks "github.com/Fi44er/sdmed/internal/module/auth/usecase/auth/mock"
	"github.com/golang/mock/gomock"
)

type MockResetPassDeps struct {
	User   *mocks.MockIUserUsecase
	Cache  *mocks.MockICache
	Notify *mocks.MockINotificationService
}

var ResetPassword = []struct {
	Name        string
	Input       string
	ExpectedErr error
	SetupMocks  func(ctx context.Context, m *MockResetPassDeps)
}{
	{
		Name:        "Success",
		Input:       "123456",
		ExpectedErr: nil,
		SetupMocks: func(ctx context.Context, m *MockResetPassDeps) {
			m.Cache.EXPECT().
				Get(ctx, "forgot_password_123456", gomock.Any()).
				DoAndReturn(func(_ context.Context, _ string, dest any) error {
					*(dest.(*string)) = "123456"
					return nil
				})
		},
	},
}
