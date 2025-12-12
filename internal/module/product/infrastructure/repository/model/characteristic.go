package product_model

import "time"

type DataType string

const (
	DataTypeString  DataType = "string"
	DataTypeNumber  DataType = "number"
	DataTypeBoolean DataType = "boolean"
	DataTypeSelect  DataType = "select"
)

type Characteristic struct {
	ID          string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Name        string    `gorm:"type:varchar(255);not null"`
	CategoryID  string    `gorm:"type:uuid;not null;index"`
	Unit        *string   `gorm:"type:varchar(50)"`                           // шт, кг, л, м и т.д.
	Description *string   `gorm:"type:text"`                                  // описание характеристики
	DataType    DataType  `gorm:"type:varchar(20);not null;default:'string'"` // string, number, boolean, enum
	IsRequired  bool      `gorm:"not null;default:false"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
	UpdatedAt   time.Time `gorm:"not null;default:now()"`

	Options  []CharOption `gorm:"foreignKey:CharacteristicID;references:ID"`
	Category Category     `gorm:"foreignKey:CategoryID;references:ID"`
}

type CharOption struct {
	ID               string    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	CharacteristicID string    `gorm:"type:uuid;not null;index"`
	Value            string    `gorm:"type:varchar(255);not null"`
	CreatedAt        time.Time `gorm:"not null;default:now()"`
}
