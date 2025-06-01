package repository

import (
	"github.com/Fi44er/sdmed/internal/module/user/entity"
	"github.com/Fi44er/sdmed/internal/module/user/infrastructure/repository/model"
)

type Converter struct{}

func (c *Converter) ToModel(entity *entity.User) *model.User {
	roles := make([]model.Role, len(entity.Roles))
	for i, r := range entity.Roles {
		roles[i] = model.Role{ID: r.ID, Name: r.Name}
	}
	return &model.User{
		ID:           entity.ID,
		Email:        entity.Email,
		PasswordHash: entity.PasswordHash,
		Name:         entity.Name,
		Surname:      entity.Surname,
		Patronymic:   entity.Patronymic,
		PhoneNumber:  entity.PhoneNumber,
		Roles:        roles,
	}
}

func (c *Converter) ToEntity(model *model.User) *entity.User {
	roles := make([]entity.Role, len(model.Roles))
	for i, r := range model.Roles {
		roles[i] = entity.Role{ID: r.ID, Name: r.Name}
	}
	return &entity.User{
		ID:           model.ID,
		Email:        model.Email,
		PasswordHash: model.PasswordHash,
		Name:         model.Name,
		Surname:      model.Surname,
		Patronymic:   model.Patronymic,
		PhoneNumber:  model.PhoneNumber,
		Roles:        roles,
	}
}
