package usecase

import (
	"context"

	"github.com/Fi44er/sdmed/internal/module/user/entity"
	"github.com/Fi44er/sdmed/internal/module/user/pkg/constant"
	"github.com/Fi44er/sdmed/pkg/logger"
	"golang.org/x/crypto/bcrypt"
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

func (u *UserUsecase) GetAll(ctx context.Context, limit, offset int) ([]entity.User, error) {
	users, err := u.repository.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		u.logger.Infof("No users found")
		return nil, constant.ErrUserNotFound
	}

	return users, nil
}

func (u *UserUsecase) GetByID(ctx context.Context, id string) (*entity.User, error) {
	user, err := u.repository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		u.logger.Infof("User with id %s not found", id)
		return nil, constant.ErrUserNotFound
	}

	return user, nil
}

func (u *UserUsecase) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user, err := u.repository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		u.logger.Infof("User with email %s not found", email)
		return nil, constant.ErrUserNotFound
	}

	return user, nil
}

func (u *UserUsecase) Create(ctx context.Context, user *entity.User) error {
	if err := u.repository.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *UserUsecase) Update(ctx context.Context, user *entity.User) error {
	if err := user.Validate(); err != nil {
		return constant.ErrInvalidUserData
	}
	if err := u.repository.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (u *UserUsecase) Delete(ctx context.Context, id string) error {
	if err := u.repository.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (u *UserUsecase) ComparePassword(user *entity.User, password string) bool {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return false
	}
	return user.ComparePassword(string(hash))
}
