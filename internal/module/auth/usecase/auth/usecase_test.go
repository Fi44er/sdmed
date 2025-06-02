package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Fi44er/sdmed/internal/config"
	"github.com/Fi44er/sdmed/internal/module/auth/entity"
	"github.com/Fi44er/sdmed/internal/module/auth/pkg/constant"
	"github.com/Fi44er/sdmed/internal/module/auth/usecase/auth"
	global_const "github.com/Fi44er/sdmed/internal/module/user/pkg/constant"
	"github.com/Fi44er/sdmed/pkg/logger"

	mocks "github.com/Fi44er/sdmed/internal/module/auth/usecase/auth/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestAuthUsecase_SignUp(t *testing.T) {
	type testCase struct {
		name        string
		input       *entity.User
		expectedErr error
	}

	tests := []testCase{
		{
			name: "Success",
			input: &entity.User{
				Email:       "new@example.com",
				PhoneNumber: "79998887766",
				Password:    "pass123",
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
		},
		{
			name: "Cache set error",
			input: &entity.User{
				Email:       "cachefail@example.com",
				PhoneNumber: "79998887766",
				Password:    "pass",
			},
			expectedErr: errors.New("cache set failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUser := mocks.NewMockIUserUsecase(ctrl)
			mockCache := mocks.NewMockICache(ctrl)
			mockNotify := mocks.NewMockINotificationService(ctrl)
			mockSession := mocks.NewMockISessionRepository(ctrl)

			cfg := &config.Config{VerifyCodeExpiredIn: 5 * time.Minute}
			logger := logger.NewLogger()
			authUC := usecase.NewAuthUsecase(logger, mockCache, cfg, mockUser, mockNotify, mockSession)

			ctx := context.Background()

			switch tt.name {
			case "Success":
				mockUser.EXPECT().
					GetByEmail(ctx, tt.input.Email).
					Return(nil, global_const.ErrUserNotFound)

				mockCache.EXPECT().
					Set(ctx, gomock.Any(), gomock.Any(), 10*time.Minute).
					Return(nil)

				mockCache.EXPECT().
					Get(ctx, gomock.Any(), gomock.Any()).
					Return(nil)

				mockCache.EXPECT().
					Set(ctx, gomock.Any(), gomock.Any(), cfg.VerifyCodeExpiredIn).
					Return(nil)

				mockNotify.EXPECT().
					Send(gomock.Any(), "smtp")

			case "User exists":
				mockUser.EXPECT().
					GetByEmail(ctx, tt.input.Email).
					Return(&entity.User{}, nil)

			case "Cache set error":
				mockUser.EXPECT().
					GetByEmail(ctx, tt.input.Email).
					Return(nil, global_const.ErrUserNotFound)

				mockCache.EXPECT().
					Set(ctx, gomock.Any(), gomock.Any(), 10*time.Minute).
					Return(errors.New("cache set failed"))
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
