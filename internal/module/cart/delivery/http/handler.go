package cart_http

import (
	"context"
	"fmt"

	cart_dto "github.com/Fi44er/sdmed/internal/module/cart/dto"
	cart_entity "github.com/Fi44er/sdmed/internal/module/cart/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	_ "github.com/Fi44er/sdmed/pkg/response"
	"github.com/Fi44er/sdmed/pkg/session"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type ICartUsecase interface {
	GetCart(ctx context.Context, userID string, cartID string) (*cart_entity.Cart, error)
	AddItem(ctx context.Context, data *cart_dto.AddCartItemRequest, userID string, cartID string) (*cart_entity.Cart, error)
	DeleteItem(ctx context.Context, itemID string, userID string, cartID string) (*cart_entity.Cart, error)
	Move(ctx context.Context, userID string, sessionCartID string) error
}

type CartHandler struct {
	usecase   ICartUsecase
	logger    *logger.Logger
	validator *validator.Validate
	converter *Converter
}

func NewCartHandler(usecase ICartUsecase, logger *logger.Logger, validator *validator.Validate) *CartHandler {
	return &CartHandler{
		usecase:   usecase,
		logger:    logger,
		validator: validator,
		converter: &Converter{},
	}
}

// Get godoc
// @Summary Get cart
// @Description Get cart for current user or session
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} cart_dto.CartResponse "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /cart [get]
func (h *CartHandler) Get(ctx *fiber.Ctx) error {
	userID := h.getUserID(ctx)
	sess := session.FromFiberContext(ctx)

	cartID := ""
	if sess != nil {
		if val, ok := sess.Get("cart_id").(string); ok {
			cartID = val
		}
	}

	cart, err := h.usecase.GetCart(ctx.Context(), userID, cartID)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   h.converter.ToCartResponse(cart),
	})
}

// AddItem godoc
// @Summary Add item to cart
// @Description Add item to cart for current user or session
// @Tags cart
// @Accept json
// @Produce json
// @Param item body cart_dto.AddCartItemRequest true "Cart Item"
// @Success 200 {object} cart_dto.CartResponse "OK"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 500 {object} response.Response "Error"
// @Router /cart/items [post]
func (h *CartHandler) AddItem(ctx *fiber.Ctx) error {
	req := new(cart_dto.AddCartItemRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	userID := h.getUserID(ctx)
	sess := session.FromFiberContext(ctx)

	fmt.Println("user id: ", userID)
	cartID := ""
	if sess != nil {
		if val, ok := sess.Get("cart_id").(string); ok {
			cartID = val
		}
	}

	cart, err := h.usecase.AddItem(ctx.Context(), req, userID, cartID)
	if err != nil {
		return err
	}

	if userID == "" && sess != nil && cart != nil {
		sess.Put("cart_id", cart.ID)
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   h.converter.ToCartResponse(cart),
	})
}

// DeleteItem godoc
// @Summary Delete item from cart
// @Description Delete item from cart for current user or session
// @Tags cart
// @Accept json
// @Produce json
// @Param id path string true "Item ID"
// @Success 200 {object} cart_dto.CartResponse "OK"
// @Failure 500 {object} response.Response "Error"
// @Router /cart/items/{id} [delete]
func (h *CartHandler) DeleteItem(ctx *fiber.Ctx) error {
	itemID := ctx.Params("id")
	userID := h.getUserID(ctx)
	sess := session.FromFiberContext(ctx)

	cartID := ""
	if sess != nil {
		if val, ok := sess.Get("cart_id").(string); ok {
			cartID = val
		}
	}

	cart, err := h.usecase.DeleteItem(ctx.Context(), itemID, userID, cartID)
	if err != nil {
		return err
	}

	return ctx.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   h.converter.ToCartResponse(cart),
	})
}

// Move godoc
// @Summary Move session cart to user
// @Description Move session cart to user after login
// @Tags cart
// @Accept json
// @Produce json
// @Param data body cart_dto.MoveCartRequest true "Move Data"
// @Success 200 {object} response.Response "OK"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 500 {object} response.Response "Error"
// @Router /cart/move [post]
func (h *CartHandler) Move(ctx *fiber.Ctx) error {
	req := new(cart_dto.MoveCartRequest)
	if err := ctx.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	sess := session.FromFiberContext(ctx)
	if sess == nil {
		return fiber.NewError(fiber.StatusBadRequest, "session not found")
	}

	sessionCartID, ok := sess.Get("cart_id").(string)
	if !ok || sessionCartID == "" {
		return ctx.Status(200).JSON(fiber.Map{"status": "success", "message": "no session cart to move"})
	}

	if err := h.usecase.Move(ctx.Context(), req.UserID, sessionCartID); err != nil {
		return err
	}

	sess.Delete("cart_id")

	return ctx.Status(200).JSON(fiber.Map{
		"status":  "success",
		"message": "cart moved successfully",
	})
}

func (h *CartHandler) getUserID(ctx *fiber.Ctx) string {
	// Implementation depends on how you store userID in context after auth
	if val, ok := ctx.Locals("user_id").(string); ok {
		return val
	}
	return ""
}
