package auth_entity

import "time"

type ActiveSession struct {
	UserID    string    `json:"user_id"`
	DeviceID  string    `json:"device_id"`
	UserRoles []string  `json:"user_roles"`
	IsShadow  bool      `json:"is_shadow"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
}

type UserSession struct {
	ID         string
	UserID     string
	UserAgent  string
	LastIP     string
	DeviceName string
	IsRevoked  bool
	// Fingerprint string

	CreatedAt  time.Time
	UpdatedAt  time.Time
	ExpiresAt  time.Time
	LastUsedAt time.Time
}

type DeviceInfo struct {
	DeviceID   string    `json:"device_id"`
	DeviceName string    `json:"device_name"`
	LastIP     string    `json:"last_ip"`
	IsCurrent  bool      `json:"is_current"`
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at"`
}
