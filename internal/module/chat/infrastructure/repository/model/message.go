package chat_models

import (
	"time"

	"gorm.io/datatypes"
)

type MessageType string

const (
	MsgTypeRegular   MessageType = "regular"
	MsgTypeSystem    MessageType = "system"
	MsgTypeOrderCard MessageType = "order_card"
	MsgTypeFeedback  MessageType = "feedback"
)

type Message struct {
	ID       string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ChatID   string `gorm:"index;not null;type:uuid"`
	SenderID string `gorm:"index;not null"`

	Type    MessageType `gorm:"type:varchar(20);default:regular"`
	Payload string      `gorm:"type:text"`

	Metadata datatypes.JSON `gorm:"type:jsonb"`

	ReplyToID *string `gorm:"type:uuid"`
	IsEdited  bool    `gorm:"default:false"`

	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
	ReadAt    *time.Time `gorm:"index"`

	// Связи
	// Attachments []Attachment `gorm:"foreignKey:MessageID"`
}

// type Attachment struct {
// 	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
// 	MessageID string `gorm:"index;not null;type:uuid"`

// 	FileName string `gorm:"type:varchar(255)"`
// 	FileSize int64
// 	MimeType string `gorm:"type:varchar(100)"`
// 	URL      string `gorm:"type:text"`

// 	CreatedAt time.Time
// }
