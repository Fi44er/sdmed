package http

import (
	"context"

	"github.com/Fi44er/sdmedik/backend/internal/module/user/dto"
	"github.com/Fi44er/sdmedik/backend/internal/module/user/entity"
	"github.com/Fi44er/sdmedik/backend/pkg/logger"
	_ "github.com/Fi44er/sdmedik/backend/pkg/response"
	"github.com/Fi44er/sdmedik/backend/pkg/session"
	"github.com/Fi44er/sdmedik/backend/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type IUserUsecase interface {
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]entity.User, error)
	Create(ctx context.Context, entity *entity.User) error
	Update(ctx context.Context, entity *entity.User) error
	Delete(ctx context.Context, id string) error
}

type UserHandler struct {
	usecase IUserUsecase

	logger    *logger.Logger
	validator *validator.Validate

	converter *Converter
}

func NewUserHandler(
	usecase IUserUsecase,
	logger *logger.Logger,
	validator *validator.Validate,
) *UserHandler {
	return &UserHandler{
		usecase:   usecase,
		logger:    logger,
		validator: validator,
		converter: &Converter{},
	}
}

// GetMy godoc
// @Summary Get my user
// @Description Get my user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} response.ResponseData{data=dto.UserResponse} "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /users/me [get]
func (h *UserHandler) GetMy(ctx *fiber.Ctx) error {
	sessFiber := session.FromFiberContext(ctx)
	h.logger.Infof("fiber session: %v", sessFiber)

	h.logger.Infof("%v", sessFiber.Get("test"))
	sessFiber.Put("test", "test")

	h.logger.Infof("%v", sessFiber.Get("test"))

	// user := ctx.Locals("user").(dto.UserResponse)
	return ctx.Status(200).JSON(fiber.Map{"status": "success", "data": "user"})
}

// GetByID godoc
// @Summary Get user by ID
// @Description Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.ResponseData{data=dto.UserResponse} "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /users/{id} [get]
func (h *UserHandler) GetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	user, err := h.usecase.GetByID(ctx.Context(), id)
	if err != nil {
		return err
	}

	response := h.converter.ToResponse(user)

	return ctx.Status(200).JSON(fiber.Map{"status": "success", "data": response})
}

// GetAll godoc
// @Summary Get all users
// @Description Get all users
// @Tags users
// @Accept json
// @Produce json
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} response.ResponseListData{data=[]dto.UserResponse} "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /users [get]
func (h *UserHandler) GetAll(ctx *fiber.Ctx) error {
	log.Info("test")
	offset := ctx.QueryInt("offset")
	limit := ctx.QueryInt("limit")

	users, err := h.usecase.GetAll(ctx.Context(), limit, offset)
	if err != nil {
		return err
	}

	response := h.converter.ToResponseSlice(users)

	return ctx.Status(200).JSON(fiber.Map{"status": "success", "data": response})
}

// Create godoc
// @Summary Create user
// @Description Create user
// @Tags users
// @Accept json
// @Produce json
// @Param user body dto.UserDTO true "User"
// @Success      200      {object}  response.Response "OK"
// @Failure      500      {object}  response.Response "Error"
// @Router /users [post]
func (h *UserHandler) Create(ctx *fiber.Ctx) error {
	dto := new(dto.UserDTO)

	entity, err := utils.ParseAndValidate(ctx, dto, h.validator, h.converter.ToEntity, h.logger)
	if err != nil {
		return err
	}

	if err := h.usecase.Create(ctx.Context(), entity); err != nil {
		h.logger.Errorf("error while create user: %s", err)
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "user created successfully",
	})
}

// Update godoc
// @Summary Update user
// @Description Update user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param user body dto.UserDTO true "User"
// @Success      200      {object}  response.Response "OK"
// @Failure      500      {object}  response.Response "Error"
// @Router /users/{id} [put]
func (h *UserHandler) Update(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	dto := new(dto.UserDTO)

	entity, err := utils.ParseAndValidate(ctx, dto, h.validator, h.converter.ToEntity, h.logger)
	if err != nil {
		return err
	}
	entity.ID = id

	if err := h.usecase.Update(ctx.Context(), entity); err != nil {
		h.logger.Errorf("error while update user: %s", err)
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "user updated successfully",
	})
}

// Delete godoc
// @Summary Delete user
// @Description Delete user
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success      200      {object}  response.Response "OK"
// @Failure      500      {object}  response.Response "Error"
// @Router /users/{id} [delete]
func (h *UserHandler) Delete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.usecase.Delete(ctx.Context(), id); err != nil {
		h.logger.Errorf("error while delete user: %s", err)
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "user deleted successfully",
	})
}
