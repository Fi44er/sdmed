package user_repository

import (
	user_entity "github.com/Fi44er/sdmed/internal/module/user/entity"
	user_model "github.com/Fi44er/sdmed/internal/module/user/infrastructure/repository/model"
)

type Converter struct{}

func (c *Converter) ToModel(entity *user_entity.User) *user_model.User {
	roles := make([]user_model.Role, len(entity.Roles))
	for i, r := range entity.Roles {
		roles[i] = user_model.Role{ID: r.ID, Name: r.Name}
	}
	return &user_model.User{
		ID:           entity.ID,
		Email:        entity.Email,
		PasswordHash: entity.PasswordHash,
		Name:         entity.Name,
		Surname:      entity.Surname,
		Patronymic:   entity.Patronymic,
		PhoneNumber:  entity.PhoneNumber,
		Roles:        roles,
		IsShadow:     entity.IsShadow,
	}
}

func (c *Converter) ToEntity(model *user_model.User) *user_entity.User {
	roles := make([]user_entity.Role, len(model.Roles))
	for i, r := range model.Roles {
		roles[i] = user_entity.Role{ID: r.ID, Name: r.Name}
	}
	return &user_entity.User{
		ID:           model.ID,
		Email:        model.Email,
		PasswordHash: model.PasswordHash,
		Name:         model.Name,
		Surname:      model.Surname,
		Patronymic:   model.Patronymic,
		PhoneNumber:  model.PhoneNumber,
		Roles:        roles,
		IsShadow:     model.IsShadow,
	}
}
