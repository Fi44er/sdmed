package characteristic_repository

import (
	"context"

	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_model "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/model"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type ICharacteristicRepository interface {
	Create(ctx context.Context, characteristic *product_entity.Characteristic) error
	CreateMany(ctx context.Context, characteristics []product_entity.Characteristic) error
	Update(ctx context.Context, characteristic *product_entity.Characteristic) error
	Delete(ctx context.Context, id string) error
	DeleteByCategory(ctx context.Context, categoryID string) error
	GetByID(ctx context.Context, id string) (*product_entity.Characteristic, error)
	GetByCategoryID(ctx context.Context, categoryID string) ([]product_entity.Characteristic, error)
	GetByCategoryAndName(ctx context.Context, categoryID, name string) (*product_entity.Characteristic, error)
}

type CharacteristicRepository struct {
	logger    *logger.Logger
	db        *gorm.DB
	converter *Converter
}

func NewCharacteristicRepository(logger *logger.Logger, db *gorm.DB) ICharacteristicRepository {
	return &CharacteristicRepository{
		logger:    logger,
		db:        db,
		converter: &Converter{},
	}
}

func (r *CharacteristicRepository) Create(ctx context.Context, characteristic *product_entity.Characteristic) error {
	r.logger.Infof("Creating characteristic: %v", characteristic.Name)

	characteristicModel := r.converter.ToModel(characteristic)
	if err := r.db.WithContext(ctx).Create(characteristicModel).Error; err != nil {
		r.logger.Errorf("Failed to create characteristic: %v", err)
		return err
	}

	characteristic.ID = characteristicModel.ID

	r.logger.Infof("Characteristic created successfully: %s (ID: %s)", characteristic.Name, characteristic.ID)
	return nil
}

func (r *CharacteristicRepository) Update(ctx context.Context, characteristic *product_entity.Characteristic) error {
	r.logger.Infof("Updating characteristic: %v", characteristic.Name)

	characteristicModel := r.converter.ToModel(characteristic)
	if err := r.db.WithContext(ctx).Model(characteristicModel).Updates(characteristicModel).Error; err != nil {
		r.logger.Errorf("Failed to update characteristic: %v", err)
		return err
	}

	r.logger.Infof("Characteristic updated successfully: %s (ID: %s)", characteristic.Name, characteristic.ID)
	return nil
}

func (r *CharacteristicRepository) Delete(ctx context.Context, id string) error {
	r.logger.Infof("Deleting characteristic with ID: %s", id)

	characteristicModel := &product_model.Characteristic{}
	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(characteristicModel).Error; err != nil {
		r.logger.Errorf("Failed to delete characteristic: %v", err)
		return err
	}

	r.logger.Infof("Characteristic deleted successfully: %s", id)
	return nil
}

func (r *CharacteristicRepository) GetByID(ctx context.Context, id string) (*product_entity.Characteristic, error) {
	r.logger.Debugf("Getting characteristic with ID: %s", id)

	characteristicModel := &product_model.Characteristic{}
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(characteristicModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("Characteristic not found: %s", id)
			return nil, nil
		}
		r.logger.Errorf("Failed to get characteristic: %v", err)
		return nil, err
	}

	characteristic := r.converter.ToEntity(characteristicModel)
	return characteristic, nil
}

func (r *CharacteristicRepository) GetByCategoryID(ctx context.Context, categoryID string) ([]product_entity.Characteristic, error) {
	r.logger.Debugf("Getting characteristics by category ID: %s", categoryID)

	characteristicModels := []*product_model.Characteristic{}
	if err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&characteristicModels).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("No characteristics found for category: %s", categoryID)
			return nil, nil
		}
		r.logger.Errorf("Failed to get characteristics: %v", err)
		return nil, err
	}

	characteristics := make([]product_entity.Characteristic, len(characteristicModels))
	for i, characteristicModel := range characteristicModels {
		characteristics[i] = *r.converter.ToEntity(characteristicModel)
	}

	return characteristics, nil
}

func (r *CharacteristicRepository) DeleteByCategory(ctx context.Context, categoryID string) error {
	r.logger.Infof("Deleting characteristics by category: %s", categoryID)

	characteristicModel := &product_model.Characteristic{}
	if err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Delete(characteristicModel).Error; err != nil {
		r.logger.Errorf("Failed to delete characteristics by category %s: %v", categoryID, err)
		return err
	}

	r.logger.Infof("Characteristics deleted successfully by category: %s", categoryID)
	return nil
}

func (r *CharacteristicRepository) CreateMany(ctx context.Context, characteristics []product_entity.Characteristic) error {
	r.logger.Infof("Creating characteristics")

	characteristicModels := make([]product_model.Characteristic, len(characteristics))
	for i, characteristic := range characteristics {
		characteristicModel := *r.converter.ToModel(&characteristic)
		characteristicModels[i] = characteristicModel
	}

	if err := r.db.WithContext(ctx).Create(characteristicModels).Error; err != nil {
		r.logger.Errorf("Failed to create characteristics: %v", err)
		return err
	}

	r.logger.Infof("Characteristics created successfully")
	return nil
}

func (r *CharacteristicRepository) GetByCategoryAndName(ctx context.Context, categoryID, name string) (*product_entity.Characteristic, error) {
	r.logger.Debugf("Getting characteristics by category: ID %s; name %s", categoryID, name)
	characteristicModel := &product_model.Characteristic{}
	if err := r.db.WithContext(ctx).Where("category_id = ? AND name = ?", categoryID, name).First(characteristicModel).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("No characteristics found for category: %s", categoryID)
			return nil, nil
		}
		r.logger.Errorf("Failed to get characteristics by category %s and name %s: %v", categoryID, name, err)
		return nil, err
	}

	entity := r.converter.ToEntity(characteristicModel)
	return entity, nil
}
