package user_model

type Permission struct {
	ID   string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	Name string `gorm:"type:varchar(255);unique;not null"`
}

func (Permission) TableName() string {
	return "user_module.permissions"
}
