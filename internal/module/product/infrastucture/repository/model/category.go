package product_model

type Category struct {
	ID   string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	Name string `gorm:"type:varchar(255);not null"`
}
