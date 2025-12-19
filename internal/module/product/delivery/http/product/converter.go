// package product_http

// import (
// 	"path"

// 	product_dto "github.com/Fi44er/sdmed/internal/module/product/dto"
// 	product_entity "github.com/Fi44er/sdmed/internal/module/product/entity"
// )

// type Converter struct{}

// func (c *Converter) ToEntityFromCreate(dto *product_dto.CreateProductRequest) *product_entity.Product {
// 	imageEntity := make([]product_entity.File, 0, len(dto.Images))
// 	charValueEntity := make([]product_entity.ProductCharValue, 0, len(dto.CharacteristicValues))
// 	for _, fileURL := range dto.Images {
// 		fileName := path.Base(fileURL)
// 		imageEntity = append(imageEntity, product_entity.File{
// 			Name: fileName,
// 		})
// 	}

// 	for _, charValue := range dto.CharacteristicValues {
// 		charValueEntity = append(charValueEntity, product_entity.ProductCharValue{
// 			CharacteristicID: charValue.CharacteristicID,
// 			StringValue:      &charValue.Value,
// 		})
// 	}

// 	return &product_entity.Product{
// 		Name:        dto.Name,
// 		Images:      imageEntity,
// 		Article:     dto.Article,
// 		Description: dto.Description,
// 		CategoryID:  dto.CategoryID,
// 		ManualPrice: &dto.ManualPrice,
// 		CharValues:  charValueEntity,
// 	}
// }

package product_http

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

// ToEntityFromCreate конвертирует DTO создания продукта в сущность продукта
func (c *Converter) ToEntityFromCreate(dto *product_dto.CreateProductRequest) *product_entity.Product {
	// Конвертация URL-адресов изображений в сущности File
	imageEntity := make([]product_entity.File, 0, len(dto.Images))
	for _, fileURL := range dto.Images {
		fileName := path.Base(fileURL)
		imageEntity = append(imageEntity, product_entity.File{
			Name: fileName,
		})
	}

	// Конвертация DTO значений характеристик в сущности
	charValuesEntity := make([]product_entity.ProductCharValue, 0, len(dto.CharacteristicValues))
	for _, cv := range dto.CharacteristicValues {
		charValuesEntity = append(charValuesEntity, product_entity.ProductCharValue{
			CharacteristicID: cv.CharacteristicID,
			StringValue:      &cv.Value,
		})
	}

	return &product_entity.Product{
		Name:        dto.Name,
		Article:     dto.Article,
		Description: dto.Description,
		ManualPrice: &dto.ManualPrice,
		IsActive:    dto.IsActive,
		Images:      imageEntity,
		CategoryID:  dto.CategoryID,
		CharValues:  charValuesEntity,
		// ID, CreatedAt, UpdatedAt устанавливаются на уровне UseCase/Repository
	}
}

// ToProductResponse конвертирует сущность продукта в DTO ответа
func (c *Converter) ToProductResponse(product *product_entity.Product) *product_dto.ProductResponse {
	if product == nil {
		return nil
	}

	// Конвертация сущностей значений характеристик в DTO
	charValuesDTO := make([]product_dto.CharValueRequest, 0, len(product.CharValues))
	for _, cv := range product.CharValues {
		// DTO CharValueRequest используется для вывода, так как структура совпадает
		charValuesDTO = append(charValuesDTO, product_dto.CharValueRequest{
			CharacteristicID: cv.CharacteristicID,
			Value:            *cv.StringValue,
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
/*
func (c *Converter) ToProductListResponse(products []product_entity.Product, count int64, page, pageSize int) *product_dto.ProductListResponse {
	// В Product DTO нет ProductListResponse, но если он появится, логика будет похожей:
	// ...
}
*/

// --- Вспомогательные функции (аналогично category_http) ---

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
	// Конфигурация для файлов
	// Предполагая, что c.config.FileLink - это "files" или аналогичный путь
	return fmt.Sprintf("%s/%s/%s", c.config.ApiUrl, c.config.FileLink, file.Name)
}
