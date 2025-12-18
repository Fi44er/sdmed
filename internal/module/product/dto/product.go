package product_dto

import "time"

type CreateProductRequest struct {
	Name                 string             `json:"name" validate:"required,min=2,max=100"`
	Article              string             `json:"article" validate:"required,min=2,max=100"`
	Description          string             `json:"description" validate:"omitempty,min=2,max=5000"` // не обязательно, но если есть - от 2 до 5000 символов
	ManualPrice          float64            `json:"manual_price" validate:"omitempty,gte=0"`         // не обязательно, но если есть - >= 0
	IsActive             bool               `json:"is_active" validate:"required"`
	Images               []string           `json:"images" validate:"required,min=1,dive,url"`
	CategoryID           string             `json:"category_id" validate:"omitempty"`
	CharacteristicValues []CharValueRequest `json:"characteristic_values" validate:"omitempty"`
}

type ProductResponse struct {
	ID                   uint64             `json:"id"`
	Name                 string             `json:"name"`
	Article              string             `json:"article"`
	Description          string             `json:"description"`
	ManualPrice          float64            `json:"manual_price"`
	IsActive             bool               `json:"is_active"`
	Images               []FileResponse     `json:"images"`
	CharacteristicValues []CharValueRequest `json:"characteristic_values"`
	CreateAt             time.Time          `json:"created_at"`
	UpdateAt             time.Time          `json:"updated_at"`
}
