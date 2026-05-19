package cart_repository

import (
	"context"

	cart_entity "github.com/Fi44er/sdmed/internal/module/cart/entity"
	cart_model "github.com/Fi44er/sdmed/internal/module/cart/infrastructure/repository/model"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
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

type CartRepository struct {
	db        *gorm.DB
	logger    *logger.Logger
	converter *Converter
}

func NewCartRepository(db *gorm.DB, logger *logger.Logger) ICartRepository {
	return &CartRepository{
		db:        db,
		logger:    logger,
		converter: &Converter{},
	}
}

func (r *CartRepository) GetByUserID(ctx context.Context, userID string) (*cart_entity.Cart, error) {
	var model cart_model.Cart
	err := r.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.converter.ToEntity(&model), nil
}

func (r *CartRepository) GetByID(ctx context.Context, id string) (*cart_entity.Cart, error) {
	var model cart_model.Cart
	err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.converter.ToEntity(&model), nil
}

func (r *CartRepository) Create(ctx context.Context, cart *cart_entity.Cart) error {
	model := r.converter.ToModel(cart)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	cart.ID = model.ID
	return nil
}

func (r *CartRepository) Update(ctx context.Context, cart *cart_entity.Cart) error {
	model := r.converter.ToModel(cart)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *CartRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&cart_model.Cart{}, "id = ?", id).Error
}

func (r *CartRepository) GetItemByCriteria(ctx context.Context, productID, cartID, iso string, options string) (*cart_entity.CartItem, error) {
	var model cart_model.CartItem
	query := r.db.WithContext(ctx).Where("product_id = ? AND cart_id = ? AND iso = ?", productID, cartID, iso)

	if options != "" && options != "[]" && options != "null" {
		query = query.Where("selected_options::text = ?", options)
	}

	err := query.First(&model).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.converter.ToItemEntity(&model), nil
}

func (r *CartRepository) CreateItem(ctx context.Context, item *cart_entity.CartItem) error {
	model := r.converter.ToItemModel(item)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	item.ID = model.ID
	return nil
}

func (r *CartRepository) UpdateItem(ctx context.Context, item *cart_entity.CartItem) error {
	model := r.converter.ToItemModel(item)
	return r.db.WithContext(ctx).Save(model).Error
}

func (r *CartRepository) DeleteItem(ctx context.Context, itemID, cartID string) error {
	return r.db.WithContext(ctx).Delete(&cart_model.CartItem{}, "id = ? AND cart_id = ?", itemID, cartID).Error
}
