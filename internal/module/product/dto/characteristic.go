package product_dto

type CreateCharacteristicRequest struct {
	Name        string   `json:"name" validate:"required,min=1,max=255"`
	Description string   `json:"description" validate:"min=1,max=255"`
	Unit        string   `json:"unit" validate:"min=1,max=255"`
	IsRequired  bool     `json:"is_required"`
	DataType    string   `json:"data_type" validate:"required"`
	Options     []string `json:"options"`
}

type CharacteristicResponse struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Unit        string       `json:"unit"`
	IsRequired  bool         `json:"is_required"`
	DataType    string       `json:"data_type"`
	Options     []CharOption `json:"options"`
}

type CharOption struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type CharValueRequest struct {
	CharacteristicID string     `json:"characteristic_id" validate:"required"`
	StringValue      string     `json:"string_value"`
	NumberValue      float64    `json:"number_value"`
	BooleanValue     bool       `json:"boolean_value"`
	OptionID         string     `json:"option_id"`
	Option           CharOption `json:"option"`
}
