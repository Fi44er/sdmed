package role_usecase

import (
	"context"

	user_entity "github.com/Fi44er/sdmed/internal/module/user/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
)

type IRoleRepository interface {
	GetAll(ctx context.Context) ([]user_entity.Role, error)
}

type IRoleUsecase interface {
	GetAll(ctx context.Context) ([]user_entity.Role, error)
}

type RoleUseCase struct {
	logger *logger.Logger
	repo   IRoleRepository
}

func NewRoleUsecase(logger *logger.Logger, repo IRoleRepository) IRoleUsecase {
	return &RoleUseCase{
		logger: logger,
		repo:   repo,
	}
}

func (u *RoleUseCase) GetAll(ctx context.Context) ([]user_entity.Role, error) {
	return u.repo.GetAll(ctx)
}
