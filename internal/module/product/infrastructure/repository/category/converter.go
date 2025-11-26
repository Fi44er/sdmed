package category_repository

import (
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	product_model "github.com/Fi44er/sdmed/internal/module/product/infrastructure/repository/model"
	"gorm.io/gorm"
)

type Converter struct{}

func (c *Converter) ToModel(entity *product_entity.Category) *product_model.Category {
	model := &product_model.Category{
		ID:        entity.ID,
		Name:      entity.Name,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}

	if entity.DeletedAt != nil {
		model.DeletedAt = gorm.DeletedAt{
			Time:  *entity.DeletedAt,
			Valid: true,
		}
	} else {
		model.DeletedAt = gorm.DeletedAt{}
	}

	return model
}

func (c *Converter) ToEntity(model *product_model.Category) *product_entity.Category {
	entity := &product_entity.Category{
		ID:   model.ID,
		Name: model.Name,
	}

	if model.DeletedAt.Valid {
		entity.DeletedAt = &model.DeletedAt.Time
	}

	return entity
}
