package entity

import (
	"fmt"
	"regexp"
)

type User struct {
	ID         string
	Name       string
	Surname    string
	Patronymic string

	Email        string
	PasswordHash string
	PhoneNumber  string

	Roles []Role
}

func (u *User) Validate() error {
	switch {
	case u.ValidatePhoneNumber(u.PhoneNumber) != nil:
		return fmt.Errorf("Invalid phone number")
	case u.ValidateEmail(u.Email) != nil:
		return fmt.Errorf("Invalid email")
	}
	return nil
}

func (u *User) AddRole(role Role) error {
	u.Roles = append(u.Roles, role)
	return nil
}

func (u *User) HasRole(role string) bool {
	for _, r := range u.Roles {
		if r.Name == role {
			return true
		}
	}
	return false
}

func (u *User) ValidatePhoneNumber(phoneNumber string) error {
	if len(phoneNumber) != 11 {
		return fmt.Errorf("Invalid phone number")
	}
	return nil
}

func (u *User) ValidateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return fmt.Errorf("Invalid email")
	}
	return nil
}

func (u *User) ComparePassword(passwordHash string) bool {
	if u.PasswordHash != passwordHash {
		return false
	}
	return true
}
