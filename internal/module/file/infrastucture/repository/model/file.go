package model

type File struct {
	ID        string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name      string `json:"name"`
	OwnerID   string `json:"owner_id"`
	OwnerType string `json:"owner_type"`
}

func (File) TableName() string {
	return "file_module.files"
}
