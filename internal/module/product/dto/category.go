package product_dto

type CategoryDTO struct {
	Name string `json:"name"`
}

type CreateCategoryDTO struct {
	Name   string   `json:"name"`
	Images []string `json:"images"`
}

type CategoryRes struct {
	ID     string    `json:"id"`
	Name   string    `json:"name"`
	Images []FileRes `json:"images"`
}

type FileRes struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}
