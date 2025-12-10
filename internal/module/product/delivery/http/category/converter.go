package category_http

import (
	"fmt"
	"math"
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
	characteristicsEntity := make([]product_entity.Characteristic, 0)
	for _, characteristic := range dto.Characteristics {
		characteristicsEntity = append(characteristicsEntity, product_entity.Characteristic{
			Name:        characteristic.Name,
			Unit:        &characteristic.Unit,
			Description: &characteristic.Description,
			DataType:    product_entity.DataType(characteristic.DataType),
			IsRequired:  characteristic.IsRequired,
		})
	}
	for _, fileURL := range dto.Images {
		fileName := path.Base(fileURL)
		imageEntity = append(imageEntity, product_entity.File{
			Name: fileName,
		})
	}
	return &product_entity.Category{
		Name:            dto.Name,
		Images:          imageEntity,
		Characteristics: characteristicsEntity,
	}
}

func (c *Converter) ToEntityFromUpdate(dto *product_dto.UpdateCategoryRequest) *product_entity.Category {
	imageEntity := make([]product_entity.File, 0)
	for _, fileURL := range dto.Images {
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

func (c *Converter) toCategoryResponses(categories []product_entity.Category, count int64, page, pageSize int) *product_dto.CategoryListResponse {
	if len(categories) == 0 {
		return &product_dto.CategoryListResponse{}
	}

	result := make([]product_dto.CategoryResponse, len(categories))
	for i, category := range categories {
		result[i] = *c.toCategoryResponse(&category)
	}
	return &product_dto.CategoryListResponse{
		Data: result,
		Pagination: product_dto.PaginationInfo{
			Total:    count,
			Page:     page,
			PageSize: pageSize,
			Pages:    int(math.Ceil(float64(count) / float64(pageSize))),
		},
	}
}

func (c *Converter) toCategoryResponse(category *product_entity.Category) *product_dto.CategoryResponse {
	if category == nil {
		return nil
	}

	return &product_dto.CategoryResponse{
		ID:              category.ID,
		Name:            category.Name,
		Slug:            category.Slug,
		Images:          c.toFileResponses(category.Images),
		Characteristics: c.toCharacteristicResponses(category.Characteristics),
	}
}

func (c *Converter) toCharacteristicResponses(characteristics []product_entity.Characteristic) []product_dto.CharacteristicResponse {
	if len(characteristics) == 0 {
		return []product_dto.CharacteristicResponse{}
	}

	characteristicResponses := make([]product_dto.CharacteristicResponse, len(characteristics))
	for i, characteristic := range characteristics {
		characteristicResponses[i] = c.toCharacteristicResponse(characteristic)
	}
	return characteristicResponses
}

func (c *Converter) toCharacteristicResponse(characteristic product_entity.Characteristic) product_dto.CharacteristicResponse {
	return product_dto.CharacteristicResponse{
		ID:          characteristic.ID,
		Name:        characteristic.Name,
		Description: *characteristic.Description,
		Unit:        *characteristic.Unit,
		DataType:    string(characteristic.DataType),
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
