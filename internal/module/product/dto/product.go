package product_dto

import "time"

type CreateProductRequest struct {
	Name                 string             `json:"name" validate:"required,min=2,max=100"`
	Article              string             `json:"article" validate:"required,min=2,max=100"`
	Description          string             `json:"description" validate:"omitempty,min=2,max=5000"`
	ManualPrice          float64            `json:"manual_price" validate:"omitempty,gte=0"`
	IsActive             bool               `json:"is_active" validate:"required"`
	Images               []string           `json:"images" validate:"required,min=1,dive,url"`
	CategoryID           string             `json:"category_id" validate:"omitempty"`
	CharacteristicValues []CharValueRequest `json:"characteristic_values" validate:"omitempty"`
}

type ProductResponse struct {
	ID                   string              `json:"id"`
	Name                 string              `json:"name"`
	Article              string              `json:"article"`
	Description          string              `json:"description"`
	ManualPrice          float64             `json:"manual_price"`
	IsActive             bool                `json:"is_active"`
	Images               []FileResponse      `json:"images"`
	CharacteristicValues []CharValueResponse `json:"characteristic_values"`
	CreateAt             time.Time           `json:"created_at"`
	UpdateAt             time.Time           `json:"updated_at"`
}

type FilterResponse struct {
	CharacteristicID   string   `json:"characteristic_id"`
	CharacteristicName string   `json:"characteristic_name"`
	DataType           string   `json:"data_type"`
	Unit               string   `json:"unit"`
	Options            []string `json:"options"`
}

type ProductQueryParams struct {
	CategoryID      string            `query:"category_id"`
	MinPrice        *float64          `query:"min_price"`
	MaxPrice        *float64          `query:"max_price"`
	Characteristics map[string]string `query:"chars"`
	Sort            string            `query:"sort"` // например: price_asc, price_desc, newest
	Page            int               `query:"page"`
	PageSize        int               `query:"page_size"`
}
