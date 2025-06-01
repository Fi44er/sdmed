package repository

import (
	"github.com/Fi44er/sdmed/internal/module/file/entity"
	"github.com/Fi44er/sdmed/internal/module/file/infrastucture/repository/model"
)

type Converter struct{}

func (c *Converter) ToModel(entity *entity.File) *model.File {
	return &model.File{
		ID:        entity.ID,
		Name:      entity.Name,
		OwnerID:   entity.OwnerID,
		OwnerType: entity.OwnerType,
	}
}

func (c *Converter) ToEntity(model *model.File) *entity.File {
	return &entity.File{
		ID:        model.ID,
		Name:      model.Name,
		OwnerID:   model.OwnerID,
		OwnerType: model.OwnerType,
	}
}
