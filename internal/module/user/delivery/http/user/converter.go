package user_http

import (
	"github.com/Fi44er/sdmed/internal/module/user/dto"
	"github.com/Fi44er/sdmed/internal/module/user/entity"
)

type Converter struct{}

func (c *Converter) ToEntity(dto *user_dto.UserDTO) *user_entity.User {
	return &user_entity.User{
		Name:         dto.Name,
		Surname:      dto.Surname,
		Patronymic:   dto.Patronymic,
		PhoneNumber:  dto.PhoneNumber,
		Email:        dto.Email,
		PasswordHash: dto.Password,
	}
}

func (c *Converter) ToResponse(entity *user_entity.User) *user_dto.UserResponse {
	roles := make([]string, len(entity.Roles))
	for i, r := range entity.Roles {
		roles[i] = r.Name
	}
	return &user_dto.UserResponse{
		ID:          entity.ID,
		Name:        entity.Name,
		Surname:     entity.Surname,
		Patronymic:  entity.Patronymic,
		PhoneNumber: entity.PhoneNumber,
		Email:       entity.Email,
		Role:        roles,
	}
}

func (c *Converter) ToResponseSlice(entities []user_entity.User) []user_dto.UserResponse {
	res := make([]user_dto.UserResponse, len(entities))
	for i, e := range entities {
		res[i] = *c.ToResponse(&e)
	}
	return res
}
