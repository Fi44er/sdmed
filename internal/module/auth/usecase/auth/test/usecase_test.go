package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/Fi44er/sdmed/internal/config"
	"github.com/Fi44er/sdmed/internal/module/auth/usecase/auth"
	"github.com/Fi44er/sdmed/pkg/logger"

	mocks "github.com/Fi44er/sdmed/internal/module/auth/usecase/auth/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthUsecase_SignUp(t *testing.T) {
	verifyCodeTTL := 5 * time.Minute
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()

			m := &mockDeps{
				user:    mocks.NewMockIUserUsecase(ctrl),
				cache:   mocks.NewMockICache(ctrl),
				notify:  mocks.NewMockINotificationService(ctrl),
				session: mocks.NewMockISessionRepository(ctrl),
			}

			cfg := &config.Config{VerifyCodeExpiredIn: verifyCodeTTL}
			log := logger.NewLogger()
			authUC := usecase.NewAuthUsecase(log, m.cache, cfg, m.user, m.notify, m.session)

			if tt.setupMocks != nil {
				tt.setupMocks(ctx, m)
			}

			err := authUC.SignUp(ctx, tt.input)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
