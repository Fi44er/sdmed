package setup_mocks

import (
	"context"
	"errors"

	mocks "github.com/Fi44er/sdmed/internal/module/auth/usecase/auth/mock"
	"github.com/golang/mock/gomock"
)

type MockValidateResetPassDeps struct {
	Cache *mocks.MockICache
}

var ValidateResetPassTests = []struct {
	Name        string
	Input       string
	ExpectedErr error
	SetupMocks  func(ctx context.Context, m *MockValidateResetPassDeps)
}{
	{
		Name:        "Success",
		Input:       "123456",
		ExpectedErr: nil,
		SetupMocks: func(ctx context.Context, m *MockValidateResetPassDeps) {
			m.Cache.EXPECT().
				Get(ctx, "forgot_password_123456", gomock.Any()).
				DoAndReturn(func(_ context.Context, _ string, dest interface{}) error {
					*(dest.(*string)) = "123456"
					return nil
				})
		},
	},
	{
		Name:        "CacheError",
		Input:       "123456",
		ExpectedErr: errors.New("cache error"),
		SetupMocks: func(ctx context.Context, m *MockValidateResetPassDeps) {
			m.Cache.EXPECT().
				Get(ctx, "forgot_password_123456", gomock.Any()).
				Return(errors.New("cache error"))
		},
	},
}
