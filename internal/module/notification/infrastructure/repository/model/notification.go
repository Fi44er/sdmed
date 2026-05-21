package model

import (
	"time"
)

type Notification struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserID    string    `gorm:"index;not null"`
	Subject   string    `gorm:"type:varchar(255);not null"`
	Content   string    `gorm:"type:text;not null"`
	IsRead    bool      `gorm:"default:false;not null"`
	CreatedAt time.Time `gorm:"not null;default:CURRENT_TIMESTAMP"`
}
