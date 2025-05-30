package http

import (
	"context"

	"github.com/Fi44er/sdmedik/backend/internal/module/auth/dto"
	"github.com/Fi44er/sdmedik/backend/internal/module/auth/entity"
	"github.com/Fi44er/sdmedik/backend/pkg/logger"
	"github.com/Fi44er/sdmedik/backend/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IAuthUsecase interface {
	// SignIn(ctx context.Context, data *dto.SignInDTO) (*dto.LoginResponse, error)
	// VerifyCode(ctx context.Context, data *dto.VerifyCodeDTO) error
	SignUp(ctx context.Context, entity *entity.User) error
	// SendCode(ctx context.Context, email string) error
	// RefreshAccessToken(ctx context.Context, data *dto.RefreshTokenDTO) (string, error)
	// SignOut(ctx context.Context, data *dto.LogoutDTO) error
}

type AuthHandler struct {
	usecase   IAuthUsecase
	validator *validator.Validate
	logger    *logger.Logger

	converter *Converter
}

func NewAuthHandler(
	usecase IAuthUsecase,
	logger *logger.Logger,
	validator *validator.Validate,
) *AuthHandler {
	return &AuthHandler{
		usecase:   usecase,
		logger:    logger,
		validator: validator,
		converter: &Converter{},
	}
}

func (h *AuthHandler) SignUp(ctx *fiber.Ctx) error {
	dto := new(dto.SignUpDTO)

	entity, err := utils.ParseAndValidate(ctx, dto, h.validator, h.converter.ToEntity, h.logger)
	if err != nil {
		return err
	}

	if err := h.usecase.SignUp(ctx.Context(), entity); err != nil {
		h.logger.Errorf("error while create user: %s", err)
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "user register successfully",
	})
}
