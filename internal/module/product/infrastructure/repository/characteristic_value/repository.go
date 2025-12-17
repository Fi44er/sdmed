package char_value_repository

import (
	"context"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_model "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/model"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type ICharValueRepository interface {
	Create(ctx context.Context, charValue *product_entity.ProductCharValue) error
	CreateMany(ctx context.Context, charValues []product_entity.ProductCharValue) error
}

type CharValueRepository struct {
	logger    *logger.Logger
	db        *gorm.DB
	converter *Converter
}

func NewCharValueRepository(logger *logger.Logger, db *gorm.DB, converter *Converter) ICharValueRepository {
	return &CharValueRepository{
		logger:    logger,
		db:        db,
		converter: converter,
	}
}

func (r *CharValueRepository) Create(ctx context.Context, charValue *product_entity.ProductCharValue) error {
	r.logger.Infof("Creating characteristic value: %v", charValue.ID)

	charValueModel := r.converter.ToModel(charValue)
	if err := r.db.WithContext(ctx).Create(charValueModel).Error; err != nil {
		r.logger.Errorf("Failed to create characteristic value '%s': %v", charValue.ID, err)
		return err
	}
	charValue.ID = charValueModel.ID

	r.logger.Infof("Characteristic value created successfully: %s (ID: %s)", charValue.ID, charValue.ID)
	return nil
}

func (r *CharValueRepository) CreateMany(ctx context.Context, charValues []product_entity.ProductCharValue) error {
	r.logger.Infof("Creating characteristic values: %v", charValues)

	charValueModels := make([]product_model.CharacteristicValue, len(charValues))
	for i, charValue := range charValues {
		charValueModels[i] = *r.converter.ToModel(&charValue)
	}

	if err := r.db.WithContext(ctx).Create(charValueModels).Error; err != nil {
		r.logger.Errorf("Failed to create characteristic values: %v", err)
		return err
	}

	r.logger.Infof("Characteristic values created successfully")
	return nil
}
