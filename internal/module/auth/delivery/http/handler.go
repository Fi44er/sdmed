package auth_http

import (
	"context"
	"time"

	"github.com/Fi44er/sdmed/internal/config"
	auth_dto "github.com/Fi44er/sdmed/internal/module/auth/dto"
	auth_entity "github.com/Fi44er/sdmed/internal/module/auth/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	_ "github.com/Fi44er/sdmed/pkg/response"
	"github.com/Fi44er/sdmed/pkg/session"
	"github.com/Fi44er/sdmed/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IAuthUsecase interface {
	SignIn(ctx context.Context, user *auth_entity.User) (*auth_entity.Tokens, error)
	VerifyCode(ctx context.Context, verifyCode *auth_entity.Code) error
	SignUp(ctx context.Context, user *auth_entity.User) error
	SendCode(ctx context.Context, sendCode *auth_entity.Code) error
	RefreshTokens(ctx context.Context, inputRefreshToken string) (*auth_entity.Tokens, error)
	SignOut(ctx context.Context) error
	SignOutAll(ctx context.Context) error // NEW
	ForgotPassword(ctx context.Context, code *auth_entity.Code) error
	ValidateResetPassword(ctx context.Context, token string) (string, error)
	ResetPassword(ctx context.Context, token string, user *auth_entity.User) error
	CreateShadowSession(ctx context.Context) (*auth_entity.Tokens, *auth_entity.User, error) // NEW
	GetUserDevices(ctx context.Context) ([]*auth_entity.DeviceInfo, error)                   // NEW
	RevokeDevice(ctx context.Context, deviceID string) error                                 // NEW
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

// @Summary CreateGuestSession
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/guest [post]
func (h *AuthHandler) CreateGuestSession(ctx *fiber.Ctx) error {
	context := h.getCtxWithSession(ctx)

	tokens, shadowUser, err := h.usecase.CreateShadowSession(context)
	if err != nil {
		h.logger.Errorf("error while creating shadow session: %s", err)
		return err
	}

	h.setCookies(ctx, tokens)

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "guest session created successfully",
		"data": fiber.Map{
			"user_id":   shadowUser.ID,
			"is_shadow": true,
		},
	})
}

// @Summary Test
// @Tags Auth
// @Accept json
// @Produce json
// @Param id query string true "ID"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/test [get]
func (h *AuthHandler) Test(ctx *fiber.Ctx) error {
	id := ctx.Query("id")
	sess := session.FromFiberContext(ctx)

	return ctx.Status(200).JSON(fiber.Map{
		"data": fiber.Map{
			"data": sess.Get(id),
		},
	})
}

// @Summary SignUp
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body auth_dto.SignUpDTO true "Sign Up"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/sign-up [post]
func (h *AuthHandler) SignUp(ctx *fiber.Ctx) error {
	dto := new(auth_dto.SignUpDTO)

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
// @Param body body auth_dto.SignInDTO true "Sign In"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/sign-in [post]
func (h *AuthHandler) SignIn(ctx *fiber.Ctx) error {
	dto := new(auth_dto.SignInDTO)

	entity, err := utils.ParseAndValidate(ctx, dto, h.validator, h.converter.ToEntitySignIn, h.logger)
	if err != nil {
		return err
	}

	context := h.getCtxWithSession(ctx)

	tokens, err := h.usecase.SignIn(context, entity)
	if err != nil {
		h.logger.Errorf("error while create user: %s", err)
		return err
	}

	h.setCookies(ctx, tokens)

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "user sign in successfully",
	})
}

// @Summary VerifyCode
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body auth_dto.VerifyCodeDTO true "Verify Code"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/verify-code [post]
func (h *AuthHandler) VerifyCode(ctx *fiber.Ctx) error {
	dto := new(auth_dto.VerifyCodeDTO)

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

	h.clearCookies(ctx)

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "sign out successfully",
	})
}

// @Summary SignOutAll
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/sign-out-all [post]
func (h *AuthHandler) SignOutAll(ctx *fiber.Ctx) error {
	context := h.getCtxWithSession(ctx)

	if err := h.usecase.SignOutAll(context); err != nil {
		h.logger.Errorf("error while sign out all: %s", err)
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "signed out from all devices except current",
	})
}

// @Summary GetDevices
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/devices [get]
func (h *AuthHandler) GetDevices(ctx *fiber.Ctx) error {
	context := h.getCtxWithSession(ctx)

	devices, err := h.usecase.GetUserDevices(context)
	if err != nil {
		h.logger.Errorf("error while getting devices: %s", err)
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   devices,
	})
}

// @Summary RevokeDevice
// @Tags Auth
// @Accept json
// @Produce json
// @Param device_id path string true "Device ID"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/devices/{device_id}/revoke [post]
func (h *AuthHandler) RevokeDevice(ctx *fiber.Ctx) error {
	deviceID := ctx.Params("device_id")
	if deviceID == "" {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": "device_id is required",
		})
	}

	context := h.getCtxWithSession(ctx)

	if err := h.usecase.RevokeDevice(context, deviceID); err != nil {
		h.logger.Errorf("error while revoking device: %s", err)
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "device revoked successfully",
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
	refreshToken := ctx.Cookies("refresh_token")
	if refreshToken == "" {
		return fiber.ErrUnauthorized
	}

	context := h.getCtxWithSession(ctx)
	tokens, err := h.usecase.RefreshTokens(context, refreshToken)
	if err != nil {
		return err
	}

	h.setCookies(ctx, tokens)

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "refresh token successfully",
	})
}

// @Summary SendCode
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body auth_dto.CodeDTO true "Send Code"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/send-code [post]
func (h *AuthHandler) SendCode(ctx *fiber.Ctx) error {
	dto := new(auth_dto.CodeDTO)

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

// @Summary ForgotPassword
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body auth_dto.CodeDTO true "Forgot Password"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(ctx *fiber.Ctx) error {
	dto := new(auth_dto.CodeDTO)
	entity, err := utils.ParseAndValidate(ctx, dto, h.validator, h.converter.ToEntityCode, h.logger)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
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
// @Param token query string true "Validate Reset password token"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/validate-reset-password [get]
func (h *AuthHandler) ValidateResetPassword(ctx *fiber.Ctx) error {
	token := ctx.Query("token")

	_, err := h.usecase.ValidateResetPassword(ctx.Context(), token)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "validate reset password token successfully",
	})
}

// @Summary ResetPassword
// @Tags Auth
// @Accept json
// @Produce json
// @Param token query string true "Reset password token"
// @Param body body auth_dto.ResetPasswordDTO true "Reset Password"
// @Success 200 {object} response.Response "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /auth/reset-password [post]
func (h *AuthHandler) ResetPassword(ctx *fiber.Ctx) error {
	token := ctx.Query("token")
	dto := new(auth_dto.ResetPasswordDTO)

	entity, err := utils.ParseAndValidate(ctx, dto, h.validator, h.converter.ToEntityResetPassword, h.logger)
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	if err := h.usecase.ResetPassword(ctx.Context(), token, entity); err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "reset password successfully",
	})
}

func (h *AuthHandler) getCtxWithSession(ctx *fiber.Ctx) context.Context {
	sess := session.FromFiberContext(ctx)

	context := context.WithValue(ctx.Context(), "session", *sess)

	return context
}

func (h *AuthHandler) setCookies(ctx *fiber.Ctx, tokens *auth_entity.Tokens) {
	ctx.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    tokens.AccessToken,
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

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/api/auth/refresh-token",
		MaxAge:   h.config.RefreshTokenMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
	})
}

func (h *AuthHandler) clearCookies(ctx *fiber.Ctx) {
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

	ctx.Cookie(&fiber.Cookie{
		Name:    "refresh_token",
		Value:   "",
		Expires: expired,
	})
}
