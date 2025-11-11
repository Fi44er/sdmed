package category_http

import (
	"fmt"

	"github.com/Fi44er/sdmed/internal/config"
	product_dto "github.com/Fi44er/sdmed/internal/module/product/dto"
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
)

type Converter struct {
	config *config.Config
}

func NewConverter(config *config.Config) *Converter {
	return &Converter{
		config: config,
	}
}

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

func (c *Converter) toCategoryResponses(categories []product_entity.Category) []product_dto.CategoryRes {
	if len(categories) == 0 {
		return []product_dto.CategoryRes{}
	}

	result := make([]product_dto.CategoryRes, len(categories))
	for i, category := range categories {
		result[i] = *c.toCategoryResponse(&category)
	}
	return result
}

func (c *Converter) toCategoryResponse(category *product_entity.Category) *product_dto.CategoryRes {
	if category == nil {
		return nil
	}

	return &product_dto.CategoryRes{
		ID:     category.ID,
		Name:   category.Name,
		Images: c.toFileResponses(category.Images),
	}
}

func (c *Converter) toFileResponses(files []product_entity.File) []product_dto.FileRes {
	if len(files) == 0 {
		return []product_dto.FileRes{}
	}

	fileResponses := make([]product_dto.FileRes, len(files))
	for i, file := range files {
		fileResponses[i] = c.toFileResponse(file)
	}
	return fileResponses
}

func (c *Converter) toFileResponse(file product_entity.File) product_dto.FileRes {
	return product_dto.FileRes{
		ID:  file.ID,
		URL: c.generateFileURL(file),
	}
}

func (c *Converter) generateFileURL(file product_entity.File) string {
	if file.ID == "" {
		return ""
	}

	return fmt.Sprintf("%s/%s/%s", c.config.ApiUrl, c.config.FileLink, file.Name)
}
