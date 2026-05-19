package cart_usecase

import (
	"context"

	cart_entity "github.com/Fi44er/sdmed/internal/module/cart/entity"
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	tru_entity "github.com/Fi44er/sdmed/internal/module/tru/entity"
)

type ICartRepository interface {
	GetByUserID(ctx context.Context, userID string) (*cart_entity.Cart, error)
	GetByID(ctx context.Context, id string) (*cart_entity.Cart, error)
	Create(ctx context.Context, cart *cart_entity.Cart) error
	Update(ctx context.Context, cart *cart_entity.Cart) error
	Delete(ctx context.Context, id string) error

	GetItemByCriteria(ctx context.Context, productID, cartID, iso string, options string) (*cart_entity.CartItem, error)
	CreateItem(ctx context.Context, item *cart_entity.CartItem) error
	UpdateItem(ctx context.Context, item *cart_entity.CartItem) error
	DeleteItem(ctx context.Context, itemID, cartID string) error
}

type IProductUsecaseAdapter interface {
	GetByID(ctx context.Context, id string) (*product_entity.Product, error)
}

type ITRUCodeUsecaseAdapter interface {
	GetByID(ctx context.Context, id string) (*tru_entity.TRUCode, error)
}

type IPromotionUsecaseAdapter interface {
	ApplyPromotions(ctx context.Context, cart *cart_entity.Cart) error
}
