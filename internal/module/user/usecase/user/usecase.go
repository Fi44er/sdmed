package usecase

import (
	"context"

	"github.com/Fi44er/sdmedik/backend/internal/module/user/entity"
	"github.com/Fi44er/sdmedik/backend/internal/module/user/pkg/constant"
	"github.com/Fi44er/sdmedik/backend/pkg/logger"
)

type IUserRepository interface {
	GetAll(ctx context.Context, limit, offset int) ([]entity.User, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, entity *entity.User) error
	Update(ctx context.Context, entity *entity.User) error
	Delete(ctx context.Context, id string) error
}

type UserUsecase struct {
	repository IUserRepository
	logger     *logger.Logger
}

func NewUserUsecase(
	repository IUserRepository,
	logger *logger.Logger,
) *UserUsecase {
	return &UserUsecase{
		repository: repository,
		logger:     logger,
	}
}

// === Query === //

func (s *UserUsecase) GetAll(ctx context.Context, limit, offset int) ([]entity.User, error) {
	users, err := s.repository.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		s.logger.Infof("No users found")
		return nil, constant.ErrUserNotFound
	}

	return users, nil
}

func (s *UserUsecase) GetByID(ctx context.Context, id string) (*entity.User, error) {
	user, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		s.logger.Infof("User with id %s not found", id)
		return nil, constant.ErrUserNotFound
	}

	return user, nil
}

func (s *UserUsecase) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := s.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		s.logger.Infof("User with email %s not found", email)
		return nil, constant.ErrUserNotFound
	}

	return user, nil
}

// === Mutation === //

func (s *UserUsecase) Create(ctx context.Context, user *entity.User) error {
	if err := user.Validate(); err != nil {
		return constant.ErrInvalidUserData
	}
	if err := s.repository.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *UserUsecase) Update(ctx context.Context, user *entity.User) error {
	if err := user.Validate(); err != nil {
		return constant.ErrInvalidUserData
	}
	if err := s.repository.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *UserUsecase) Delete(ctx context.Context, id string) error {
	if err := s.repository.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}
