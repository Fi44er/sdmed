package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Fi44er/sdmed/internal/config"
	"github.com/Fi44er/sdmed/internal/module/auth/usecase/auth"
	"github.com/Fi44er/sdmed/pkg/logger"

	mocks "github.com/Fi44er/sdmed/internal/module/auth/usecase/auth/mock"
	"github.com/Fi44er/sdmed/internal/module/auth/usecase/auth/setup_mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthUsecase_SignUp(t *testing.T) {
	verifyCodeTTL := 5 * time.Minute
	for _, tt := range setup_mocks.SignUpTests {
		t.Run(tt.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			m := &setup_mocks.MockSignUpDeps{
				User:   mocks.NewMockIUserUsecase(ctrl),
				Cache:  mocks.NewMockICache(ctrl),
				Notify: mocks.NewMockINotificationService(ctrl),
			}

			cfg := &config.Config{VerifyCodeExpiredIn: verifyCodeTTL}
			log := logger.NewLogger()
			authUC := usecase.NewAuthUsecase(log, m.Cache, cfg, m.User, m.Notify, nil, nil)

			if tt.SetupMocks != nil {
				tt.SetupMocks(ctx, m)
			}

			err := authUC.SignUp(ctx, tt.Input)

			if tt.ExpectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.ExpectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthUsecase_SignIn(t *testing.T) {
	for _, tt := range setup_mocks.SignInTests {
		t.Run(tt.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			m := &setup_mocks.MockSignInDeps{
				User:         mocks.NewMockIUserUsecase(ctrl),
				Session:      mocks.NewMockISessionRepository(ctrl),
				TokenService: mocks.NewMockITokenService(ctrl),
			}

			cfg := &config.Config{
				AccessTokenPrivateKey:  setup_mocks.ACCESS_TOKEN_PRIVATE_KEY,
				AccessTokenExpiresIn:   time.Hour,
				RefreshTokenPrivateKey: setup_mocks.REFRESH_TOKEN_PRIVATE_KEY,
				RefreshTokenExpiresIn:  time.Hour,
			}
			log := logger.NewLogger()
			authUC := usecase.NewAuthUsecase(log, nil, cfg, m.User, nil, m.Session, m.TokenService)

			if tt.SetupMocks != nil {
				tt.SetupMocks(ctx, m)
			}

			_, err := authUC.SignIn(ctx, tt.Input)

			if tt.ExpectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.ExpectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthUsecase_SignOut(t *testing.T) {
	for _, tt := range setup_mocks.SignOutTests {
		t.Run(tt.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			m := &setup_mocks.SignOutDeps{
				Session: mocks.NewMockISessionRepository(ctrl),
			}

			log := logger.NewLogger()
			authUC := usecase.NewAuthUsecase(log, nil, nil, nil, nil, m.Session, nil)

			if tt.SetupMocks != nil {
				tt.SetupMocks(ctx, m)
			}

			err := authUC.SignOut(ctx)

			if tt.ExpectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.ExpectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthUsecase_RefreshToken(t *testing.T) {
	for _, tt := range setup_mocks.RefreshTokenTests {
		t.Run(tt.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			m := &setup_mocks.MockRefreshTokenDeps{
				User:         mocks.NewMockIUserUsecase(ctrl),
				Session:      mocks.NewMockISessionRepository(ctrl),
				TokenService: mocks.NewMockITokenService(ctrl),
			}

			log := logger.NewLogger()
			cfg := &config.Config{
				AccessTokenPrivateKey: setup_mocks.ACCESS_TOKEN_PRIVATE_KEY,
				AccessTokenExpiresIn:  time.Hour,
				RefreshTokenPublicKey: setup_mocks.REFRESH_TOKEN_PUBLIC_KEY,
			}
			authUC := usecase.NewAuthUsecase(log, nil, cfg, m.User, nil, m.Session, m.TokenService)

			if tt.SetupMocks != nil {
				tt.SetupMocks(ctx, m)
			}

			_, err := authUC.RefreshAccessToken(ctx)

			if tt.ExpectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.ExpectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAuthUsecase_VerifyCode(t *testing.T) {
	for _, tt := range setup_mocks.VerifyCodeTests {
		t.Run(tt.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			m := &setup_mocks.MockVerifyCodeDeps{
				User:  mocks.NewMockIUserUsecase(ctrl),
				Cache: mocks.NewMockICache(ctrl),
			}

			log := logger.NewLogger()
			authUC := usecase.NewAuthUsecase(log, m.Cache, nil, m.User, nil, nil, nil)

			if tt.SetupMocks != nil {
				tt.SetupMocks(ctx, m)
			}

			err := authUC.VerifyCode(ctx, tt.Input)

			if tt.ExpectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.ExpectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
