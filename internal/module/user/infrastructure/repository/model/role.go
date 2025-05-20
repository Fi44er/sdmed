package model

type Role struct {
	ID          string       `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	Name        string       `gorm:"type:varchar(255);unique;not null"`
	Permissions []Permission `gorm:"many2many:user_module.role_permissions;"`
}

func (Role) TableName() string {
	return "user_module.roles"
}
