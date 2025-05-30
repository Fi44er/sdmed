package http

import (
	"github.com/Fi44er/sdmedik/backend/internal/module/auth/dto"
	"github.com/Fi44er/sdmedik/backend/internal/module/auth/entity"
)

type Converter struct{}

func (c *Converter) ToEntity(dto *dto.SignUpDTO) *entity.User {
	return &entity.User{
		Email:       dto.Email,
		Password:    dto.Password,
		FIO:         dto.FIO,
		PhoneNumber: dto.PhoneNumber,
	}
}
