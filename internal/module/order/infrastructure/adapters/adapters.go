package order_adapters

import (
	"context"
	"fmt"

	cart_entity "github.com/Fi44er/sdmed/internal/module/cart/entity"
	cart_usecase "github.com/Fi44er/sdmed/internal/module/cart/usecase/cart"
	notification_service "github.com/Fi44er/sdmed/internal/module/notification/service"
	order_entity "github.com/Fi44er/sdmed/internal/module/order/entity"
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_usecase "github.com/Fi44er/sdmed/internal/module/product/usecase/product"
	tru_usecase "github.com/Fi44er/sdmed/internal/module/tru/usecase"
)

// CartAdapter
type CartAdapter struct {
	cartUsecase *cart_usecase.CartUsecase
}

func NewCartAdapter(cartUsecase *cart_usecase.CartUsecase) *CartAdapter {
	return &CartAdapter{cartUsecase: cartUsecase}
}

func (a *CartAdapter) GetCart(ctx context.Context, userID string) (*cart_entity.Cart, error) {
	return a.cartUsecase.GetCart(ctx, userID, "")
}

func (a *CartAdapter) ClearCart(ctx context.Context, cartID string) error {
	return a.cartUsecase.Clear(ctx, cartID)
}

// ProductAdapter
type ProductAdapter struct {
	productUsecase product_usecase.IProductUsecase
}

func NewProductAdapter(productUsecase product_usecase.IProductUsecase) *ProductAdapter {
	return &ProductAdapter{productUsecase: productUsecase}
}

func (a *ProductAdapter) GetByID(ctx context.Context, id string) (*product_entity.Product, error) {
	return a.productUsecase.GetByID(ctx, id)
}

// TRUAdapter
type TRUAdapter struct {
	truUsecase tru_usecase.ITRUCodeUseCase
}

func NewTRUAdapter(truUsecase tru_usecase.ITRUCodeUseCase) *TRUAdapter {
	return &TRUAdapter{truUsecase: truUsecase}
}

func (a *TRUAdapter) GetByCode(ctx context.Context, code string) (string, error) {
	res, err := a.truUsecase.GetByCode(ctx, code)
	if err != nil || res == nil {
		return "", err
	}
	return res.Code, nil
}

// NotificationAdapter
type NotificationAdapter struct {
	notificationService *notification_service.NotificationService
}

func NewNotificationAdapter(notificationService *notification_service.NotificationService) *NotificationAdapter {
	return &NotificationAdapter{notificationService: notificationService}
}

func (a *NotificationAdapter) SendOrderNotification(ctx context.Context, order *order_entity.Order) {
	msg := &notification_service.Message{
		Recipient:    "admin@sdmedik.ru", // Multiple recipients could be added
		Subject:      "Новый заказ",
		TemplatePath: "./templates/order.html",
		Data:         order,
	}
	a.notificationService.Send(msg, "smtp")
}

// ChatAdapter (Placeholder)
type ChatAdapter struct{}

func NewChatAdapter() *ChatAdapter {
	return &ChatAdapter{}
}

func (a *ChatAdapter) AddEndMsgID(ctx context.Context, chatID string) (string, error) {
	return "", fmt.Errorf("chat module not yet implemented")
}
