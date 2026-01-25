package auth_entity

import "time"

type ActiveSession struct {
	UserID      string   `json:"u_id"`
	IsShadow    bool     `json:"is_g"`
	Roles       []string `json:"roles"`
	RefreshHash string   `json:"r_hash"` // Нужен для эндпоинта /refresh
	// Fingerprint string   `json:"fp"`
}

type UserSession struct {
	// ID — это наш DeviceID (UUID). Он зашит в Access и Refresh токены.
	ID     string
	UserID string

	// Храним хеш Refresh-токена. Сам токен в базу НЕ КЛАДЕМ.
	RefreshHash string

	// Метаданные для отображения в списке устройств
	UserAgent string
	LastIP    string
	// Fingerprint string

	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}
