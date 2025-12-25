package auth_entity

type User struct {
	ID          string
	Email       string
	Password    string
	PhoneNumber string
	FIO         string
	Roles       []Role
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
