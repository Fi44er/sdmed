package product_dto

import "time"

// CREATE - строгая валидация
type CreateCategoryRequest struct {
	Name   string   `json:"name" validate:"required,min=1,max=255"`
	Images []string `json:"images" validate:"required,min=1,dive,url"`
}

// UPDATE - частичное обновление
type UpdateCategoryRequest struct {
	ID     string    `json:"id" validate:"min=1,max=255"`
	Name   *string   `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Images *[]string `json:"images,omitempty" validate:"omitempty,dive,url"`
}

// RESPONSE - полная информация
type CategoryResponse struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Images    []FileResponse `json:"images"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// RESPONSE - краткая информация для списков
type CategoryShortResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image,omitempty"` // URL первого изображения
}

// RESPONSE - для файлов
type FileResponse struct {
	ID   string `json:"id"`
	URL  string `json:"url"`
	Name string `json:"name,omitempty"`
}

// RESPONSE - список категорий с пагинацией
type CategoryListResponse struct {
	Data       []CategoryShortResponse `json:"data"`
	Pagination PaginationInfo          `json:"pagination"`
}

type PaginationInfo struct {
	Total    int `json:"total"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Pages    int `json:"pages"`
}
