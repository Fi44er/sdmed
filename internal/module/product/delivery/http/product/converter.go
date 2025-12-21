package product_http

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

func (c *Converter) ToEntityFromCreate(dto *product_dto.CreateProductRequest) *product_entity.Product {
	imageEntity := make([]product_entity.File, 0, len(dto.Images))
	for _, fileURL := range dto.Images {
		fileName := path.Base(fileURL)
		imageEntity = append(imageEntity, product_entity.File{
			Name: fileName,
		})
	}

	charValuesEntity := make([]product_entity.ProductCharValue, 0, len(dto.CharacteristicValues))
	for _, cv := range dto.CharacteristicValues {
		charValuesEntity = append(charValuesEntity, product_entity.ProductCharValue{
			CharacteristicID: cv.CharacteristicID,
			StringValue:      &cv.Value,
		})
	}

	categoryID := &dto.CategoryID
	if *categoryID == "" {
		categoryID = nil
	}

	return &product_entity.Product{
		Name:        dto.Name,
		Article:     dto.Article,
		Description: dto.Description,
		ManualPrice: &dto.ManualPrice,
		IsActive:    dto.IsActive,
		Images:      imageEntity,
		CategoryID:  categoryID,
		CharValues:  charValuesEntity,
	}
}

func (c *Converter) ToProductResponse(product *product_entity.Product) *product_dto.ProductResponse {
	if product == nil {
		return nil
	}

	charValuesDTO := make([]product_dto.CharValueResponse, 0, len(product.CharValues))
	for _, cv := range product.CharValues {
		charValuesDTO = append(charValuesDTO, product_dto.CharValueResponse{
			CharacteristicID:   cv.CharacteristicID,
			Value:              cv.GetStringValue(),
			CharacteristicName: cv.CharacteristicName,
		})
	}

	return &product_dto.ProductResponse{
		ID:                   product.ID,
		Name:                 product.Name,
		Article:              product.Article,
		Description:          product.Description,
		ManualPrice:          *product.ManualPrice,
		IsActive:             product.IsActive,
		Images:               c.toFileResponses(product.Images),
		CharacteristicValues: charValuesDTO,
		CreateAt:             product.CreatedAt,
		UpdateAt:             product.UpdatedAt,
	}
}

// ToProductListResponse конвертирует список сущностей продукта в DTO ответа (если нужен список)
func (c *Converter) ToProductListResponse(products []product_entity.Product, count int64, page, pageSize int) *dto_utils.ListResponse[product_dto.ProductResponse] {
	if len(products) == 0 {
		return &dto_utils.ListResponse[product_dto.ProductResponse]{}
	}

	result := make([]product_dto.ProductResponse, len(products))
	for i, product := range products {
		result[i] = *c.ToProductResponse(&product)
	}

	return &dto_utils.ListResponse[product_dto.ProductResponse]{
		Data: result,
		Pagination: dto_utils.PaginationInfo{
			Total:    count,
			Page:     page,
			PageSize: pageSize,
			Pages:    int(math.Ceil(float64(count) / float64(pageSize))),
		},
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

func (c *Converter) ToFilterResponses(filters []product_entity.Filter) []product_dto.FilterResponse {
	if len(filters) == 0 {
		return []product_dto.FilterResponse{}
	}

	filterResponses := make([]product_dto.FilterResponse, len(filters))
	for i, filter := range filters {
		filterResponses[i] = c.toFilterResponse(filter)
	}
	return filterResponses
}

func (c *Converter) toFilterResponse(filter product_entity.Filter) product_dto.FilterResponse {
	return product_dto.FilterResponse{
		CharacteristicID:   filter.CharacteristicID,
		CharacteristicName: filter.CharacteristicName,
		DataType:           filter.DataType,
		Unit:               filter.Unit,
		Options:            filter.Options,
	}
}

func (c *Converter) ToFilterEntity(filter product_dto.ProductQueryParams) product_entity.ProductFilterParams {
	return product_entity.ProductFilterParams{
		Page:            filter.Page,
		PageSize:        filter.PageSize,
		CategoryID:      filter.CategoryID,
		MinPrice:        filter.MinPrice,
		MaxPrice:        filter.MaxPrice,
		Sort:            filter.Sort,
		Characteristics: filter.Characteristics,
	}
}
