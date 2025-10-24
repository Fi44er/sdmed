package model

type TRUCodePrice struct {
	ID        string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	TRUCodeID string  `gorm:"type:uuid;not null;"`
	RegionID  string  `gorm:"type:uuid;not null;"`
	Price     float64 `gorm:"type:float;not null"`
}

func (TRUCodePrice) TableName() string {
	return "tru_module.tru_code_prices"
}
