package chat_models

import (
	"time"

	"github.com/lib/pq" // Используется для работы с массивами строк в Postgres
	// Нужен для работы с JSON (Metadata)
	"gorm.io/gorm"
)

// --- Enums ---

type ChatStatus string

const (
	StatusNew      ChatStatus = "new"
	StatusOpen     ChatStatus = "open"
	StatusPending  ChatStatus = "pending"
	StatusClosed   ChatStatus = "closed"
	StatusArchived ChatStatus = "archived"
)

type Chat struct {
	ID         string  `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	ClientID   string  `gorm:"index;not null"`
	OperatorID *string `gorm:"index"`

	Subject  string     `gorm:"type:varchar(255)"`
	Status   ChatStatus `gorm:"type:varchar(20);default:new;index"`
	Priority int        `gorm:"default:1"`

	Tags pq.StringArray `gorm:"type:text[]"`

	CreatedAt time.Time `gorm:"index"`
	UpdatedAt time.Time
	ClosedAt  *time.Time     `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Duration в Go обычно хранится как int64 (наносекунды) в базе
	FirstResponseTime *int64
	Rating            int `gorm:"default:0"`

	// Связи
	Messages     []Message         `gorm:"foreignKey:ChatID"`
	Participants []ChatParticipant `gorm:"foreignKey:ChatID"`
}

type ChatParticipant struct {
	ID     uint   `gorm:"primaryKey"`
	ChatID string `gorm:"uniqueIndex:idx_chat_user;not null;type:uuid"`
	UserID string `gorm:"uniqueIndex:idx_chat_user;not null"`

	Role     string    `gorm:"type:varchar(20);default:customer"`
	JoinedAt time.Time `gorm:"autoCreateTime"`
}
