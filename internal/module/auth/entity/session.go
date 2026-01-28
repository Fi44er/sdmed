package auth_entity

import "time"

type ActiveSession struct {
	UserID      string   `json:"u_id"`
	UserRoles   []string `json:"roles"`
	DeviceID    string   `json:"device_id"`
	DeviceName  string   `json:"device_name"`
	IP          string   `json:"ip"`
	RefreshHash string   `json:"r_hash"`
	IsShadow    bool     `json:"is_g"`
	// Fingerprint string   `json:"fp"`
}

type UserSession struct {
	ID          string
	UserID      string
	RefreshHash string
	UserAgent   string
	LastIP      string
	DeviceName  string
	IsRevoked   bool
	// Fingerprint string

	CreatedAt time.Time
	UpdatedAt time.Time
	ExpiresAt time.Time
}

type DeviceInfo struct {
	DeviceID   string    `json:"device_id"`
	DeviceName string    `json:"device_name"`
	UserAgent  string    `json:"user_agent"`
	LastIP     string    `json:"last_ip"`
	IsCurrent  bool      `json:"is_current"`
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at"`
}
