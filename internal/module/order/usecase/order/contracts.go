package order_usecase

import (
	"context"

	cart_entity "github.com/Fi44er/sdmed/internal/module/cart/entity"
	order_entity "github.com/Fi44er/sdmed/internal/module/order/entity"
	order_payment "github.com/Fi44er/sdmed/internal/module/order/infrastructure/payment"
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
)

type IOrderRepository interface {
	Create(ctx context.Context, order *order_entity.Order) error
	Update(ctx context.Context, order *order_entity.Order) error
	GetByID(ctx context.Context, id string) (*order_entity.Order, error)
	GetByUserID(ctx context.Context, userID string) ([]order_entity.Order, error)
	GetAll(ctx context.Context, offset, limit int) ([]order_entity.Order, error)
	AddItems(ctx context.Context, items []order_entity.OrderItem) error
}

type IPaymentProvider interface {
	CreateInvoice(ctx context.Context, email, phone, fio string, amount float64, items []order_payment.PayKeeperItem) (string, error)
}

type ICartAdapter interface {
	GetCart(ctx context.Context, userID string) (*cart_entity.Cart, error)
	ClearCart(ctx context.Context, cartID string) error
}

type IProductAdapter interface {
	GetByID(ctx context.Context, id string) (*product_entity.Product, error)
}

type ITRUAdapter interface {
	GetByCode(ctx context.Context, code string) (string, error)
}

type INotificationAdapter interface {
	SendOrderNotification(ctx context.Context, order *order_entity.Order)
}

type IChatAdapter interface {
	AddEndMsgID(ctx context.Context, chatID string) (string, error)
}
