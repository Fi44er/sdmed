package repository

import (
	file_entity "github.com/Fi44er/sdmed/internal/module/file/entity"
	"github.com/Fi44er/sdmed/internal/module/file/infrastucture/repository/model"
)

type Converter struct{}

func (c *Converter) ToModel(entity *file_entity.File) *model.File {
	return &model.File{
		ID:        entity.ID,
		Name:      entity.Name,
		OwnerID:   entity.OwnerID,
		OwnerType: entity.OwnerType,
	}
}

func (c *Converter) ToEntity(model *model.File) *file_entity.File {
	return &file_entity.File{
		ID:        model.ID,
		Name:      model.Name,
		OwnerID:   model.OwnerID,
		OwnerType: model.OwnerType,
	}
}
