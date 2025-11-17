package product_dto

type CategoryDTO struct {
	Name string `json:"name"`
}

type UpdateCatgoryDTO struct {
	Name   string   `json:"name" validate:"required,min=1,max=255"`
	Images []string `json:"images" validate:"dive,url"`
}

type CreateCategoryDTO struct {
	Name   string   `json:"name" validate:"required,min=1,max=255"`
	Images []string `json:"images" validate:"dive,url"`
}

type CategoryResponse struct {
	ID     string         `json:"id"`
	Name   string         `json:"name"`
	Images []FileResponse `json:"images"`
}

type FileResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}
