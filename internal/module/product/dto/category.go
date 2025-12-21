package product_dto

import "time"

type CreateCategoryRequest struct {
	Name            string                        `json:"name" validate:"required,min=1,max=255"`
	Images          []string                      `json:"images" validate:"required,min=1,dive,url"`
	Characteristics []CreateCharacteristicRequest `json:"characteristics" validate:"dive"`
}

type UpdateCategoryRequest struct {
	ID              string                        `json:"id" validate:"min=1,max=255"`
	Name            *string                       `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Images          []string                      `json:"images,omitempty" validate:"omitempty,dive,url"`
	Characteristics []CreateCharacteristicRequest `json:"characteristics,omitempty" validate:"omitempty,dive"`
}

type CategoryResponse struct {
	ID              string                   `json:"id"`
	Name            string                   `json:"name"`
	Slug            string                   `json:"slug"`
	Images          []FileResponse           `json:"images"`
	Characteristics []CharacteristicResponse `json:"characteristics"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

type CategoryShortResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image,omitempty"`
}

type FileResponse struct {
	ID   string `json:"id"`
	URL  string `json:"url"`
	Name string `json:"name,omitempty"`
}
