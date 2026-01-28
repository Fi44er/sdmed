package auth_entity

import "time"

type User struct {
	ID              string
	Email           string
	Password        string
	PhoneNumber     string
	FIO             string
	IsShadow        bool
	ShadowCreatedAt *time.Time
	ShadowExpiresAt *time.Time
	Roles           []Role
}

type Role struct {
	ID          string
	Name        string
	Permissions []Permission
}

type Permission struct {
	ID   string
	Name string
}
