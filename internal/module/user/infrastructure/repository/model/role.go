package user_model

type Role struct {
	ID          string       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	Name        string       `gorm:"type:varchar(255);unique;not null"`
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}
