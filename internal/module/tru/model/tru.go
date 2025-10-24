package model

type TRUCode struct {
	ID       string         `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	Code     string         `gorm:"type:varchar(30);not null;unique"`
	IsCustom bool           `gorm:"type:bool;not null;default:false"`
	Prices   []TRUCodePrice `gorm:"foreignKey:TRUCodeID"`
}

func (TRUCode) TableName() string {
	return "tru_module.tru_codes"
}
