package role_repository

import (
	"context"

	user_entity "github.com/Fi44er/sdmed/internal/module/user/entity"
	user_model "github.com/Fi44er/sdmed/internal/module/user/infrastructure/repository/model"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type IRoleRepository interface {
	GetAll(ctx context.Context) ([]user_entity.Role, error)
}

type RoleRepository struct {
	logger   *logger.Logger
	db       *gorm.DB
	converte *Converter
}

func NewRoleRepository(logger *logger.Logger, db *gorm.DB) *RoleRepository {
	return &RoleRepository{
		logger:   logger,
		db:       db,
		converte: &Converter{},
	}
}

func (r *RoleRepository) GetAll(ctx context.Context) ([]user_entity.Role, error) {
	var roles []user_model.Role
	if err := r.db.WithContext(ctx).Preload("Permissions").Find(&roles).Error; err != nil {
		return nil, err
	}

	entities := make([]user_entity.Role, len(roles))
	for i, role := range roles {
		entities[i] = *r.converte.ToEntity(&role)
	}

	return entities, nil
}
