package usecase

import (
	"context"

	"github.com/Fi44er/sdmed/internal/module/product/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
)

type ICategoryRepository interface {
	Create(ctx context.Context, entity *product_entity.Category) error
	GetByID(ctx context.Context, id string) (*product_entity.Category, error)
}

type CategoryUsecase struct {
	repository ICategoryRepository
	logger     *logger.Logger
}
