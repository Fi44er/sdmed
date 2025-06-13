package http

import (
	"context"
	"time"

	"github.com/Fi44er/sdmed/internal/config"
	"github.com/Fi44er/sdmed/internal/module/auth/dto"
	"github.com/Fi44er/sdmed/internal/module/auth/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	_ "github.com/Fi44er/sdmed/pkg/response"
	"github.com/Fi44er/sdmed/pkg/session"
	"github.com/Fi44er/sdmed/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IAuthUsecase interface {
	SignIn(ctx context.Context, user *entity.User) (*entity.Tokens, error)
	VerifyCode(ctx context.Context, verifyCode *entity.Code) error
	SignUp(ctx context.Context, user *entity.User) error
	SendCode(ctx context.Context, sendCode *entity.Code) error
	RefreshAccessToken(ctx context.Context) (string, error)
	SignOut(ctx context.Context) error
	ForgotPassword(ctx context.Context, code *entity.Code) error
	ValidateResetPassword(ctx context.Context, token string) (string, error)
	ResetPassword(ctx context.Context, token string, user *entity.User) error
}

type AuthHandler struct {
	usecase   IAuthUsecase
	validator *validator.Validate
	logger    *logger.Logger
	config    *config.Config

	converter *Converter
}

func NewAuthHandler(
	usecase IAuthUsecase,
	logger *logger.Logger,
	validator *validator.Validate,
	config *config.Config,
) *AuthHandler {
	return &AuthHandler{
		usecase:   usecase,
		logger:    logger,
		validator: validator,
		converter: &Converter{},
		config:    config,
	}
}

// @Summary SignUp
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.SignUpDTO true "Sign Up"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/sign-up [post]
func (h *AuthHandler) SignUp(ctx *fiber.Ctx) error {
	dto := new(dto.SignUpDTO)

	entity, err := utils.ParseAndValidate(ctx, dto, h.validator, h.converter.ToEntitySignUp, h.logger)
	if err != nil {
		return err
	}

	if err := h.usecase.SignUp(ctx.Context(), entity); err != nil {
		h.logger.Errorf("error while create user: %s", err)
		return err
	}

	return ctx.Status(201).JSON(fiber.Map{
		"status":  "success",
		"message": "user sign up successfully",
	})
}

// @Summary SignIn
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.SignInDTO true "Sign In"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/sign-in [post]
func (h *AuthHandler) SignIn(ctx *fiber.Ctx) error {
	dto := new(dto.SignInDTO)

	entity, err := utils.ParseAndValidate(ctx, dto, h.validator, h.converter.ToEntitySignIn, h.logger)
	if err != nil {
		return err
	}

	context := h.getCtxWithSession(ctx)

	accessToken, err := h.usecase.SignIn(context, entity)
	if err != nil {
		h.logger.Errorf("error while create user: %s", err)
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken.AccessToken,
		Path:     "/",
		MaxAge:   h.config.AccessTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   h.config.AccessTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: false,
	})

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "user sign in successfully",
	})
}

// @Summary VerifyCode
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.VerifyCodeDTO true "Verify Code"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/verify-code [post]
func (h *AuthHandler) VerifyCode(ctx *fiber.Ctx) error {
	dto := new(dto.VerifyCodeDTO)

	entity, err := utils.ParseAndValidate(ctx, dto, h.validator, h.converter.ToEntityVerifyCode, h.logger)
	if err != nil {
		return err
	}

	if err := h.usecase.VerifyCode(ctx.Context(), entity); err != nil {
		h.logger.Errorf("error while create user: %s", err)
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "verify code successfully",
	})
}

// @Summary SignOut
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/sign-out [post]
func (h *AuthHandler) SignOut(ctx *fiber.Ctx) error {
	context := h.getCtxWithSession(ctx)

	if err := h.usecase.SignOut(context); err != nil {
		h.logger.Errorf("error while sign out: %s", err)
		return err
	}

	expired := time.Now().Add(-time.Hour * 24)

	ctx.Cookie(&fiber.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: expired,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:    "logged_in",
		Value:   "",
		Expires: expired,
	})

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "sign out successfully",
	})
}

// @Summary RefreshToken
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/refresh-token [post]
func (h *AuthHandler) RefreshToken(ctx *fiber.Ctx) error {
	context := h.getCtxWithSession(ctx)
	accessToken, err := h.usecase.RefreshAccessToken(context)
	if err != nil {
		return err
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   h.config.AccessTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
	})

	ctx.Cookie(&fiber.Cookie{
		Name:     "logged_in",
		Value:    "true",
		Path:     "/",
		MaxAge:   h.config.AccessTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: false,
	})

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "refresh token successfully",
	})
}

// @Summary SendCode
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.CodeDTO true "Send Code"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/send-code [post]
func (h *AuthHandler) SendCode(ctx *fiber.Ctx) error {
	dto := new(dto.CodeDTO)

	entity, err := utils.ParseAndValidate(ctx, dto, h.validator, h.converter.ToEntityCode, h.logger)
	if err != nil {
		return err
	}

	if err := h.usecase.SendCode(ctx.Context(), entity); err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "send code successfully",
	})
}

func (h *AuthHandler) getCtxWithSession(ctx *fiber.Ctx) context.Context {
	sess := session.FromFiberContext(ctx)
	context := context.WithValue(ctx.Context(), "session", *sess)

	return context
}

// @Summary ForgotPassword
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body dto.CodeDTO true "Forgot Password"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(ctx *fiber.Ctx) error {
	dto := new(dto.CodeDTO)
	entity, err := utils.ParseAndValidate(ctx, dto, h.validator, h.converter.ToEntityCode, h.logger)
	if err != nil {
		return err
	}

	if err := h.usecase.ForgotPassword(ctx.Context(), entity); err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "send resset link successfully",
	})
}

// @Summary ValidateResetPassword
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/validate-reset-password [get]
func (h *AuthHandler) ValidateResetPassword(ctx *fiber.Ctx) error {
	token := ctx.Query("token")

	_, err := h.usecase.ValidateResetPassword(ctx.Context(), token)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "validate reset password token successfully",
	})
}

func (h *AuthHandler) ResetPassword(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	dto := new(dto.ResetPasswordDTO)

	entity, err := utils.ParseAndValidate(ctx, dto, h.validator, h.converter.ToEntityResetPassword, h.logger)
	if err != nil {
		return err
	}

	if err := h.usecase.ResetPassword(ctx.Context(), token, entity); err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "reset password successfully",
	})
}
