package product_model

type Characteristic struct {
	ID         string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	Name       string  `gorm:"type:varchar(255);not null"`
	CategoryID string  `gorm:"type:varchar(255);not null"`
	Unit       *string `gorm:"type:varchar(255);"`
}

func (Characteristic) TableName() string {
	return "product_module.characteristics"
}
