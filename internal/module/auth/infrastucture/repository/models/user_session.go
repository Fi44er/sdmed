package auth_models

import (
	"time"
)

type UserSession struct {
	ID          string    `gorm:"primaryKey;type:varchar(255)"` // device_id from JWT
	UserID      string    `gorm:"index:idx_user_sessions_user_id;type:varchar(255);not null"`
	RefreshHash string    `gorm:"type:varchar(255);not null"`
	UserAgent   string    `gorm:"type:text"`
	LastIP      string    `gorm:"type:varchar(45)"`
	DeviceName  string    `gorm:"type:varchar(255)"`
	IsRevoked   bool      `gorm:"default:false;index:idx_user_sessions_revoked"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	ExpiresAt   time.Time `gorm:"index:idx_user_sessions_expires_at;not null"`
	LastUsedAt  time.Time `gorm:"index:idx_user_sessions_last_used_at;not null"`
	// User User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}
