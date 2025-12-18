package product_http

import (
	"path"

	product_dto "github.com/Fi44er/sdmed/internal/module/product/dto"
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
)

type Converter struct{}

func (c *Converter) ToEntityFromCreate(dto *product_dto.CreateProductRequest) *product_entity.Product {
	imageEntity := make([]product_entity.File, 0, len(dto.Images))
	charValueEntity := make([]product_entity.ProductCharValue, 0, len(dto.CharacteristicValues))
	for _, fileURL := range dto.Images {
		fileName := path.Base(fileURL)
		imageEntity = append(imageEntity, product_entity.File{
			Name: fileName,
		})
	}

	for _, charValue := range dto.CharacteristicValues {
		charValueEntity = append(charValueEntity, product_entity.ProductCharValue{
			CharacteristicID: charValue.CharacteristicID,
			StringValue:      &charValue.Value,
		})
	}

	return &product_entity.Product{
		Name:        dto.Name,
		Images:      imageEntity,
		Article:     dto.Article,
		Description: dto.Description,
		CategoryID:  dto.CategoryID,
		ManualPrice: &dto.ManualPrice,
		CharValues:  charValueEntity,
	}
}
