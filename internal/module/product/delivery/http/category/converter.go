package category_http

import (
	"fmt"
	"math"
	"path"

	"github.com/Fi44er/sdmed/internal/config"
	product_dto "github.com/Fi44er/sdmed/internal/module/product/dto"
	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
	dto_utils "github.com/Fi44er/sdmed/pkg/utils/dto"
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
		options := make([]product_entity.CharOption, len(characteristic.Options))
		for i, option := range characteristic.Options {
			options[i] = product_entity.CharOption{
				Value: option,
			}
		}
		characteristicsEntity = append(characteristicsEntity, product_entity.Characteristic{
			Name:        characteristic.Name,
			Unit:        &characteristic.Unit,
			Description: &characteristic.Description,
			DataType:    product_entity.DataType(characteristic.DataType),
			IsRequired:  characteristic.IsRequired,
			Options:     options,
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

func (c *Converter) ToCategoryListResponse(categories []product_entity.Category, count int64, page, pageSize int) *dto_utils.ListResponse[product_dto.CategoryResponse] {
	if len(categories) == 0 {
		return &dto_utils.ListResponse[product_dto.CategoryResponse]{}
	}

	result := make([]product_dto.CategoryResponse, len(categories))
	for i, category := range categories {
		result[i] = *c.toCategoryResponse(&category)
	}
	return &dto_utils.ListResponse[product_dto.CategoryResponse]{
		Data: result,
		Pagination: dto_utils.PaginationInfo{
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
		CreatedAt:       category.CreatedAt,
		UpdatedAt:       category.UpdatedAt,
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

	options := make([]product_dto.CharOption, len(characteristic.Options))
	for i, option := range characteristic.Options {
		options[i] = product_dto.CharOption{
			ID:    option.ID,
			Value: option.Value,
		}
	}

	return product_dto.CharacteristicResponse{
		ID:          characteristic.ID,
		Name:        characteristic.Name,
		Description: *characteristic.Description,
		Unit:        *characteristic.Unit,
		DataType:    string(characteristic.DataType),
		Options:     options,
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
