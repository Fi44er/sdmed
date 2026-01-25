package user_entity

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID         string
	Name       string
	Surname    string
	Patronymic string

	Email        string
	PasswordHash string
	PhoneNumber  string

	IsShadow bool

	Roles []Role
}

func (u *User) Validate() error {
	switch {
	case u.ValidatePhoneNumber() != nil:
		return fmt.Errorf("Invalid phone number")
	case u.ValidateEmail() != nil:
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

func (u *User) ValidatePhoneNumber() error {
	if len(u.PhoneNumber) != 11 {
		return fmt.Errorf("Invalid phone number")
	}
	return nil
}

func (u *User) ValidateEmail() error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return fmt.Errorf("Invalid email")
	}
	return nil
}

func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}
