package setup_mocks

import (
	"context"

	"github.com/Fi44er/sdmed/internal/module/auth/pkg/constant"
	mocks "github.com/Fi44er/sdmed/internal/module/auth/usecase/auth/mock"
)

type SignOutDeps struct {
	Session *mocks.MockISessionRepository
}

var SignOutTests = []struct {
	Name        string
	SetupMocks  func(context.Context, *SignOutDeps)
	ExpectedErr error
}{
	{
		Name: "Success",
		SetupMocks: func(ctx context.Context, m *SignOutDeps) {
			m.Session.EXPECT().
				DeleteSessionInfo(ctx).
				Return(nil)
		},
	},
	{
		Name:        "SessionDeleteError",
		ExpectedErr: constant.ErrSessionInfoNotFound,
		SetupMocks: func(ctx context.Context, m *SignOutDeps) {
			m.Session.EXPECT().
				DeleteSessionInfo(ctx).
				Return(constant.ErrSessionInfoNotFound)
		},
	},
}
