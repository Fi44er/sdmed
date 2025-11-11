package category_http

import (
	product_dto "github.com/Fi44er/sdmed/internal/module/product/dto"
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
)

type Converter struct{}

func (c *Converter) ToEntity(dto *product_dto.CreateCategoryDTO) *product_entity.Category {
	imageEntity := make([]product_entity.File, 0)
	for _, imageName := range dto.Images {
		imageEntity = append(imageEntity, product_entity.File{
			Name: imageName,
		})
	}
	return &product_entity.Category{
		Name:   dto.Name,
		Images: imageEntity,
	}
}
