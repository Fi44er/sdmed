package setup_mocks

import (
	"context"
	mocks "github.com/Fi44er/sdmed/internal/module/auth/usecase/auth/mock"
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
}{}
