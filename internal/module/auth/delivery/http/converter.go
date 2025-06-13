package http

import (
	"github.com/Fi44er/sdmed/internal/module/auth/dto"
	"github.com/Fi44er/sdmed/internal/module/auth/entity"
)

type Converter struct{}

func (c *Converter) ToEntitySignUp(dto *dto.SignUpDTO) *entity.User {
	return &entity.User{
		Email:       dto.Email,
		Password:    dto.Password,
		FIO:         dto.FIO,
		PhoneNumber: dto.PhoneNumber,
	}
}

func (c *Converter) ToEntitySignIn(dto *dto.SignInDTO) *entity.User {
	return &entity.User{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func (c *Converter) ToEntityVerifyCode(dto *dto.VerifyCodeDTO) *entity.Code {
	return &entity.Code{
		Email: dto.Email,
		Code:  dto.Code,
	}
}

func (c *Converter) ToEntityCode(dto *dto.CodeDTO) *entity.Code {
	return &entity.Code{
		Email: dto.Email,
	}
}

func (c *Converter) ToEntityResetPassword(dto *dto.ResetPasswordDTO) *entity.User {
	return &entity.User{
		Password: dto.Password,
	}
}
