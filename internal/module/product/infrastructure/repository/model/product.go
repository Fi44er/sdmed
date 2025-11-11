package product_model

type Product struct {
	ID              string           `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	Article         string           `gorm:"type:varchar(255);not null"`
	Name            string           `gorm:"type:varchar(255);not null"`
	Description     string           `gorm:"type:varchar(255);"`
	Price           float64          `gorm:"type:float;not null"`
	CategoryID      string           `gorm:"type:varchar(255);not null"`
	Characteristics []Characteristic `gorm:"many2many:product_characteristics;"`
}

func (Product) TableName() string {
	return "product_module.products"
}
