package http

import (
	"context"

	"github.com/Fi44er/sdmedik/backend/internal/module/auth/dto"
)

type IAuthUsecase interface {
	SignIn(ctx context.Context, data *dto.LoginDTO) (*dto.LoginResponse, error)
	VerifyCode(ctx context.Context, data *dto.VerifyCodeDTO) error
	SignUp(ctx context.Context, data *dto.RegisterDTO) error
	SendCode(ctx context.Context, email string) error
	RefreshAccessToken(ctx context.Context, data *dto.RefreshTokenDTO) (string, error)
	SignOut(ctx context.Context, data *dto.LogoutDTO) error
}
