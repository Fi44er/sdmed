package auth_http

import (
	"github.com/Fi44er/sdmed/internal/module/auth/dto"
	"github.com/Fi44er/sdmed/internal/module/auth/entity"
)

type Converter struct{}

func (c *Converter) ToEntitySignUp(dto *auth_dto.SignUpDTO) *auth_entity.User {
	return &auth_entity.User{
		Email:       dto.Email,
		Password:    dto.Password,
		FIO:         dto.FIO,
		PhoneNumber: dto.PhoneNumber,
	}
}

func (c *Converter) ToEntitySignIn(dto *auth_dto.SignInDTO) *auth_entity.User {
	return &auth_entity.User{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func (c *Converter) ToEntityVerifyCode(dto *auth_dto.VerifyCodeDTO) *auth_entity.Code {
	return &auth_entity.Code{
		Email: dto.Email,
		Code:  dto.Code,
	}
}

func (c *Converter) ToEntityCode(dto *auth_dto.CodeDTO) *auth_entity.Code {
	return &auth_entity.Code{
		Email: dto.Email,
	}
}

func (c *Converter) ToEntityResetPassword(dto *auth_dto.ResetPasswordDTO) *auth_entity.User {
	return &auth_entity.User{
		Password: dto.Password,
	}
}
