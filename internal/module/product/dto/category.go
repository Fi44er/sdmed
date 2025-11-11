package product_dto

type CategoryDTO struct {
	Name string `json:"name"`
}

type CreateCategoryDTO struct {
	Name   string   `json:"name"`
	Images []string `json:"images"`
}
