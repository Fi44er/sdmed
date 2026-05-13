package tru_model

type TRUCodePrice struct {
	ID        string  `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	TRUCodeID string  `gorm:"type:uuid;not null;"`
	RegionIso string  `gorm:"not null;"`
	Price     float64 `gorm:"type:float;not null"`
}
