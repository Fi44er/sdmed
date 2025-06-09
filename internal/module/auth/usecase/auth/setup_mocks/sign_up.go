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

var tmpCodeTTL = 10 * time.Minute

type MockSignUpDeps struct {
	User   *mocks.MockIUserUsecase
	Cache  *mocks.MockICache
	Notify *mocks.MockINotificationService
}

var SignUpTests = []struct {
	Name        string
	Input       *entity.User
	ExpectedErr error
	SetupMocks  func(ctx context.Context, m *MockSignUpDeps)
}{
	{
		Name: "Success",
		Input: &entity.User{
			Email:       "new@example.com",
			PhoneNumber: "79998887766",
			Password:    "pass123",
		},
		SetupMocks: func(ctx context.Context, m *MockSignUpDeps) {
			gomock.InOrder(
				m.User.EXPECT().
					GetByEmail(ctx, "new@example.com").
					Return(nil, constant.ErrUserNotFound),

				m.Cache.EXPECT().
					Set(ctx, gomock.Any(), gomock.Any(), tmpCodeTTL).
					Return(nil),

				m.Cache.EXPECT().
					Get(ctx, gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, key string, dest interface{}) error {
						if userPtr, ok := dest.(*entity.User); ok {
							*userPtr = entity.User{
								Email:       "new@example.com",
								PhoneNumber: "79998887766",
								Password:    "hashed_pass",
							}
						}
						return nil
					}),

				m.Cache.EXPECT().
					Set(ctx, gomock.Any(), gomock.Any(), 5*time.Minute).
					Return(nil),

				m.Notify.EXPECT().
					Send(gomock.Any(), "smtp"),
			)
		},
	},
	{
		Name: "InvalidPhone",
		Input: &entity.User{
			Email:       "bad@example.com",
			PhoneNumber: "123",
			Password:    "pass",
		},
		ExpectedErr: constant.ErrInvalidPhoneNumber,
	},
	{
		Name: "UserExists",
		Input: &entity.User{
			Email:       "exists@example.com",
			PhoneNumber: "79998887766",
			Password:    "pass",
		},
		ExpectedErr: constant.ErrUserAlreadyExists,
		SetupMocks: func(ctx context.Context, m *MockSignUpDeps) {
			m.User.EXPECT().
				GetByEmail(ctx, "exists@example.com").
				Return(&entity.User{}, nil)
		},
	},
	{
		Name: "CacheSetError",
		Input: &entity.User{
			Email:       "cachefail@example.com",
			PhoneNumber: "79998887766",
			Password:    "pass",
		},
		ExpectedErr: errors.New("cache set failed"),
		SetupMocks: func(ctx context.Context, m *MockSignUpDeps) {
			m.User.EXPECT().
				GetByEmail(ctx, "cachefail@example.com").
				Return(nil, constant.ErrUserNotFound)

			m.Cache.EXPECT().
				Set(ctx, gomock.Any(), gomock.Any(), tmpCodeTTL).
				Return(errors.New("cache set failed"))
		},
	},
}
