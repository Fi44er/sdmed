package order_http

import (
	"strconv"

	order_dto "github.com/Fi44er/sdmed/internal/module/order/dto"
	_ "github.com/Fi44er/sdmed/internal/module/order/entity"
	order_usecase "github.com/Fi44er/sdmed/internal/module/order/usecase/order"
	"github.com/Fi44er/sdmed/pkg/logger"
	_ "github.com/Fi44er/sdmed/pkg/response"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	usecase   *order_usecase.OrderUsecase
	logger    *logger.Logger
	validator *validator.Validate
}

func NewOrderHandler(usecase *order_usecase.OrderUsecase, logger *logger.Logger, validator *validator.Validate) *OrderHandler {
	return &OrderHandler{
		usecase:   usecase,
		logger:    logger,
		validator: validator,
	}
}

// CreateOrder godoc
// @Summary Create order from cart
// @Description Create a new order using items in the user's current cart and generate a PayKeeper payment link
// @Tags orders
// @Accept json
// @Produce json
// @Param order body order_dto.CreateOrderRequest true "Order Details"
// @Success 200 {object} map[string]string "OK"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /orders [post]
func (h *OrderHandler) CreateOrder(ctx *fiber.Ctx) error {
	var req order_dto.CreateOrderRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	userID, _ := ctx.Locals("user_id").(string)

	link, err := h.usecase.CreateOrder(ctx.Context(), &req, userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"payment_link": link})
}

// CreateDirectOrder godoc
// @Summary Create direct order
// @Description Create a direct order for a single product without using the cart and generate a PayKeeper payment link
// @Tags orders
// @Accept json
// @Produce json
// @Param order body order_dto.CreateDirectOrderRequest true "Direct Order Details"
// @Success 200 {object} map[string]string "OK"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /orders/direct [post]
func (h *OrderHandler) CreateDirectOrder(ctx *fiber.Ctx) error {
	var req order_dto.CreateDirectOrderRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	link, err := h.usecase.CreateDirectOrder(ctx.Context(), &req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"payment_link": link})
}

// GetMyOrders godoc
// @Summary Get my orders
// @Description Retrieve a list of all orders belonging to the authenticated user
// @Tags orders
// @Accept json
// @Produce json
// @Success 200 {array} order_entity.Order "OK"
// @Failure 401 {object} response.Response "Unauthorized"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /orders/my [get]
func (h *OrderHandler) GetMyOrders(ctx *fiber.Ctx) error {
	userID, ok := ctx.Locals("user_id").(string)
	if !ok || userID == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	orders, err := h.usecase.GetMyOrders(ctx.Context(), userID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(orders)
}

// GetAllOrders godoc
// @Summary Get all orders (Admin only)
// @Description Retrieve a paginated list of all orders in the system
// @Tags orders
// @Accept json
// @Produce json
// @Param offset query int false "Offset" default(0)
// @Param limit query int false "Limit" default(10)
// @Success 200 {array} order_entity.Order "OK"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /orders/admin [get]
func (h *OrderHandler) GetAllOrders(ctx *fiber.Ctx) error {
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))

	orders, err := h.usecase.GetAllOrders(ctx.Context(), offset, limit)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(orders)
}

// ChangeStatus godoc
// @Summary Change order status (Admin only)
// @Description Update the status of an existing order
// @Tags orders
// @Accept json
// @Produce json
// @Param status body order_dto.ChangeOrderStatusRequest true "Status Details"
// @Success 200 "OK"
// @Failure 400 {object} response.Response "Bad Request"
// @Failure 500 {object} response.Response "Internal Server Error"
// @Router /orders/admin/status [post]
func (h *OrderHandler) ChangeStatus(ctx *fiber.Ctx) error {
	var req order_dto.ChangeOrderStatusRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := h.usecase.ChangeStatus(ctx.Context(), &req); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *OrderHandler) RegisterRoutes(router fiber.Router, authMiddleware, adminMiddleware fiber.Handler) {
	group := router.Group("/orders")

	group.Post("/", h.CreateOrder)
	group.Post("/direct", h.CreateDirectOrder)
	group.Get("/my", authMiddleware, h.GetMyOrders)

	// Admin routes
	adminGroup := group.Group("/admin", adminMiddleware)
	adminGroup.Get("/", h.GetAllOrders)
	adminGroup.Post("/status", h.ChangeStatus)
}
