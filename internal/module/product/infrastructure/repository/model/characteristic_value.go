package product_model

type CharacteristicValue struct {
	ID               string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	CharacteristicID string `gorm:"type:varchar(255);not null"`
	Value            string `gorm:"type:varchar(255);not null"`
	ProductID        string `gorm:"type:varchar(255);not null"`
}

func (CharacteristicValue) TableName() string {
	return "product_module.characteristic_values"
}
