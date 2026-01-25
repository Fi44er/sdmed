package auth_models

import (
	"time"

	"gorm.io/gorm"
)

type UserSession struct {
	ID     string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID string `gorm:"index;not null;type:uuid"`

	RefreshHash string `gorm:"not null;index"`

	UserAgent string `gorm:"type:text"`
	LastIP    string `gorm:"type:varchar(45)"`
	// Fingerprint string `gorm:"index"` // Для дополнительной защиты

	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time `gorm:"index"`

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
