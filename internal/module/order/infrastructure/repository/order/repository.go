package order_repository

import (
	"context"

	order_entity "github.com/Fi44er/sdmed/internal/module/order/entity"
	order_model "github.com/Fi44er/sdmed/internal/module/order/infrastructure/repository/model"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type IOrderRepository interface {
	Create(ctx context.Context, order *order_entity.Order) error
	Update(ctx context.Context, order *order_entity.Order) error
	GetByID(ctx context.Context, id string) (*order_entity.Order, error)
	GetByUserID(ctx context.Context, userID string) ([]order_entity.Order, error)
	GetAll(ctx context.Context, offset, limit int) ([]order_entity.Order, error)
	AddItems(ctx context.Context, items []order_entity.OrderItem) error
}

type OrderRepository struct {
	db        *gorm.DB
	logger    *logger.Logger
	converter *Converter
}

func NewOrderRepository(db *gorm.DB, logger *logger.Logger) IOrderRepository {
	return &OrderRepository{
		db:        db,
		logger:    logger,
		converter: &Converter{},
	}
}

func (r *OrderRepository) Create(ctx context.Context, order *order_entity.Order) error {
	model := r.converter.ToModel(order)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		return err
	}
	order.ID = model.ID
	return nil
}

func (r *OrderRepository) Update(ctx context.Context, order *order_entity.Order) error {
	model := r.converter.ToModel(order)
	return r.db.WithContext(ctx).Model(&order_model.Order{}).Where("id = ?", order.ID).Updates(model).Error
}

func (r *OrderRepository) GetByID(ctx context.Context, id string) (*order_entity.Order, error) {
	var model order_model.Order
	if err := r.db.WithContext(ctx).Preload("Items").Where("id = ?", id).First(&model).Error; err != nil {
		return nil, err
	}
	return r.converter.ToEntity(&model), nil
}

func (r *OrderRepository) GetByUserID(ctx context.Context, userID string) ([]order_entity.Order, error) {
	var models []order_model.Order
	if err := r.db.WithContext(ctx).Preload("Items").Where("user_id = ?", userID).Order("created_at desc").Find(&models).Error; err != nil {
		return nil, err
	}

	entities := make([]order_entity.Order, len(models))
	for i, m := range models {
		entities[i] = *r.converter.ToEntity(&m)
	}
	return entities, nil
}

func (r *OrderRepository) GetAll(ctx context.Context, offset, limit int) ([]order_entity.Order, error) {
	var models []order_model.Order
	if err := r.db.WithContext(ctx).Preload("Items").Offset(offset).Limit(limit).Order("created_at desc").Find(&models).Error; err != nil {
		return nil, err
	}

	entities := make([]order_entity.Order, len(models))
	for i, m := range models {
		entities[i] = *r.converter.ToEntity(&m)
	}
	return entities, nil
}

func (r *OrderRepository) AddItems(ctx context.Context, items []order_entity.OrderItem) error {
	models := make([]order_model.OrderItem, len(items))
	for i, item := range items {
		models[i] = *r.converter.ToItemModel(&item)
	}
	return r.db.WithContext(ctx).Create(&models).Error
}
