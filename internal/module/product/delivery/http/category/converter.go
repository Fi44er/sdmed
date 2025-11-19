package category_http

import (
	"fmt"
	"path"

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

func (c *Converter) ToEntityFromCreate(dto *product_dto.CreateCategoryRequest) *product_entity.Category {
	imageEntity := make([]product_entity.File, 0)
	for _, fileURL := range dto.Images {
		fileName := path.Base(fileURL)
		imageEntity = append(imageEntity, product_entity.File{
			Name: fileName,
		})
	}
	return &product_entity.Category{
		Name:   dto.Name,
		Images: imageEntity,
	}
}

func (c *Converter) ToEntityFromUpdate(dto *product_dto.UpdateCategoryRequest) *product_entity.Category {
	imageEntity := make([]product_entity.File, 0)
	for _, fileURL := range *dto.Images {
		fileName := path.Base(fileURL)
		imageEntity = append(imageEntity, product_entity.File{
			Name: fileName,
		})
	}
	return &product_entity.Category{
		ID:     dto.ID,
		Name:   *dto.Name,
		Images: imageEntity,
	}
}

func (c *Converter) toCategoryResponses(categories []product_entity.Category) []product_dto.CategoryResponse {
	if len(categories) == 0 {
		return []product_dto.CategoryResponse{}
	}

	result := make([]product_dto.CategoryResponse, len(categories))
	for i, category := range categories {
		result[i] = *c.toCategoryResponse(&category)
	}
	return result
}

func (c *Converter) toCategoryResponse(category *product_entity.Category) *product_dto.CategoryResponse {
	if category == nil {
		return nil
	}

	return &product_dto.CategoryResponse{
		ID:     category.ID,
		Name:   category.Name,
		Images: c.toFileResponses(category.Images),
	}
}

func (c *Converter) toFileResponses(files []product_entity.File) []product_dto.FileResponse {
	if len(files) == 0 {
		return []product_dto.FileResponse{}
	}

	fileResponses := make([]product_dto.FileResponse, len(files))
	for i, file := range files {
		fileResponses[i] = c.toFileResponse(file)
	}
	return fileResponses
}

func (c *Converter) toFileResponse(file product_entity.File) product_dto.FileResponse {
	return product_dto.FileResponse{
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
