package product_model

import "time"

// type Characteristic struct {
// 	ID         string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
// 	Name       string  `gorm:"type:varchar(255);not null"`
// 	CategoryID string  `gorm:"type:varchar(255);not null"`
// 	Unit       *string `gorm:"type:varchar(255);"` // шт кг и тд..
// }

type DataType string

const (
	DataTypeString  DataType = "string"
	DataTypeNumber  DataType = "number"
	DataTypeBoolean DataType = "boolean"
	DataTypeEnum    DataType = "enum"
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

	Category Category `gorm:"foreignKey:CategoryID;references:ID"`
	// ProductCharacteristics []ProductCharacteristic `gorm:"foreignKey:CharacteristicID"` // если будет связь с товарами
}

func (Characteristic) TableName() string {
	return "product_module.characteristics"
}
