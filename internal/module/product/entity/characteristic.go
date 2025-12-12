package product_entity

import (
	"time"

	product_constant "github.com/Fi44er/sdmed/internal/module/product/pkg"
)

type DataType string

const (
	DataTypeString  DataType = "string"
	DataTypeNumber  DataType = "number"
	DataTypeBoolean DataType = "boolean"
	DataTypeSelect  DataType = "select"
)

type Characteristic struct {
	ID          string
	Name        string
	CategoryID  string
	Description *string
	Unit        *string
	DataType    DataType
	IsRequired  bool
	Options     []CharOption
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CharOption struct {
	ID               string
	CharacteristicID string
	Value            string
	CreatedAt        time.Time
}

func (e *Characteristic) ValidateDataType() error {
	switch e.DataType {
	case DataTypeString, DataTypeNumber, DataTypeBoolean, DataTypeSelect:
		return nil
	default:
		return product_constant.ErrInvalidDataTypeCharacteristic
	}
}
