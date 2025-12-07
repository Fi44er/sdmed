package product_dto

type CreateCharacteristicRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"min=1,max=255"`
	Unit        string `json:"unit" validate:"min=1,max=255"`
	IsRequired  bool   `json:"is_required"`
	DataType    string `json:"data_type" validate:"required"`
}

type CharacteristicResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Unit        string `json:"unit"`
	IsRequired  bool   `json:"is_required"`
	DataType    string `json:"data_type"`
}
