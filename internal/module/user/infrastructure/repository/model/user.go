package user_model

import "time"

type User struct {
	ID           string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4();"`
	Email        string `gorm:"type:varchar(100);unique;not null"`
	PasswordHash string `gorm:"type:varchar(255);not null"`
	Name         string `gorm:"type:varchar(255);not null"`
	Surname      string `gorm:"type:varchar(255);"`
	Patronymic   string `gorm:"type:varchar(255);"`
	PhoneNumber  string `gorm:"type:varchar(255);"`
	IsShadow     bool   `gorm:"default:false"`
	Roles        []Role `gorm:"many2many:user_roles;"`

	ShadowCreatedAt *time.Time `gorm:"index:idx_users_shadow"`
	ShadowExpiresAt *time.Time `gorm:"index:idx_users_shadow"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
