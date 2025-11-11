package model

import "time"

type File struct {
	ID        string     `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string     `json:"name"`
	OwnerID   *string    `json:"owner_id"`
	OwnerType *string    `json:"owner_type"`
	ExpiresAt *time.Time `json:"expires_at"`
	CreatedAt time.Time  `json:"created_at"`
	Status    string     `gorm:"column:status;type:file_status;not null" json:"status"`
}

func (File) TableName() string {
	return "file_module.files"
}
