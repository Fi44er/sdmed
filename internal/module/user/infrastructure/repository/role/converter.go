package role_repository

import (
	user_entity "github.com/Fi44er/sdmed/internal/module/user/entity"
	user_model "github.com/Fi44er/sdmed/internal/module/user/infrastructure/repository/model"
)

type Converter struct{}

func (c *Converter) ToModel(entity *user_entity.Role) *user_model.Role {
	permissionsModels := make([]user_model.Permission, len(entity.Permissions))
	for i, permission := range entity.Permissions {
		permissionsModels[i] = user_model.Permission{
			ID:   permission.ID,
			Name: permission.Name,
		}
	}
	return &user_model.Role{
		ID:          entity.ID,
		Name:        entity.Name,
		Permissions: permissionsModels,
	}
}

func (c *Converter) ToEntity(model *user_model.Role) *user_entity.Role {
	permissionsEntities := make([]user_entity.Permission, len(model.Permissions))
	for i, permission := range model.Permissions {
		permissionsEntities[i] = user_entity.Permission{
			ID:   permission.ID,
			Name: permission.Name,
		}
	}
	return &user_entity.Role{
		ID:          model.ID,
		Name:        model.Name,
		Permissions: permissionsEntities,
	}
}
