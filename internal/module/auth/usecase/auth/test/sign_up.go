package usecase_test

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

type mockDeps struct {
	user    *mocks.MockIUserUsecase
	cache   *mocks.MockICache
	notify  *mocks.MockINotificationService
	session *mocks.MockISessionRepository
}

var tests = []struct {
	name        string
	input       *entity.User
	expectedErr error
	setupMocks  func(ctx context.Context, m *mockDeps)
}{
	{
		name: "Success",
		input: &entity.User{
			Email:       "new@example.com",
			PhoneNumber: "79998887766",
			Password:    "pass123",
		},
		setupMocks: func(ctx context.Context, m *mockDeps) {
			gomock.InOrder(
				m.user.EXPECT().
					GetByEmail(ctx, "new@example.com").
					Return(nil, constant.ErrUserNotFound),

				m.cache.EXPECT().
					Set(ctx, gomock.Any(), gomock.Any(), 10*time.Minute).
					Return(nil),

				m.cache.EXPECT().
					Get(ctx, gomock.Any(), gomock.Any()).
					Return(nil),

				m.cache.EXPECT().
					Set(ctx, gomock.Any(), gomock.Any(), 5*time.Minute).
					Return(nil),

				m.notify.EXPECT().
					Send(gomock.Any(), "smtp"),
			)
		},
	},
	{
		name: "Invalid phone",
		input: &entity.User{
			Email:       "bad@example.com",
			PhoneNumber: "123",
			Password:    "pass",
		},
		expectedErr: constant.ErrInvalidPhoneNumber,
	},
	{
		name: "User exists",
		input: &entity.User{
			Email:       "exists@example.com",
			PhoneNumber: "79998887766",
			Password:    "pass",
		},
		expectedErr: constant.ErrUserAlreadyExists,
		setupMocks: func(ctx context.Context, m *mockDeps) {
			m.user.EXPECT().
				GetByEmail(ctx, "exists@example.com").
				Return(&entity.User{}, nil)
		},
	},
	{
		name: "Cache set error",
		input: &entity.User{
			Email:       "cachefail@example.com",
			PhoneNumber: "79998887766",
			Password:    "pass",
		},
		expectedErr: errors.New("cache set failed"),
		setupMocks: func(ctx context.Context, m *mockDeps) {
			m.user.EXPECT().
				GetByEmail(ctx, "cachefail@example.com").
				Return(nil, constant.ErrUserNotFound)

			m.cache.EXPECT().
				Set(ctx, gomock.Any(), gomock.Any(), tmpCodeTTL).
				Return(errors.New("cache set failed"))
		},
	},
}
