package adapters

import (
	"context"

	authEntity "github.com/Fi44er/sdmed/internal/module/auth/entity"
	constantAuth "github.com/Fi44er/sdmed/internal/module/auth/pkg/constant"
	userEntity "github.com/Fi44er/sdmed/internal/module/user/entity"
	constantUser "github.com/Fi44er/sdmed/internal/module/user/pkg/constant"
	userUsecase "github.com/Fi44er/sdmed/internal/module/user/usecase/user"
)

type UserUsecaseAdapter struct {
	userUsecase *userUsecase.UserUsecase
}

func NewUserUsecaseAdapter(userUsecase *userUsecase.UserUsecase) *UserUsecaseAdapter {
	return &UserUsecaseAdapter{
		userUsecase: userUsecase,
	}
}

func (a *UserUsecaseAdapter) GetByEmail(ctx context.Context, email string) (*authEntity.User, error) {
	user, err := a.userUsecase.GetByEmail(ctx, email)
	if err != nil {
		if err == constantUser.ErrUserNotFound {
			return nil, constantAuth.ErrUserNotFound
		}
		return nil, err
	}

	return toAuthUser(user), nil
}

func (a *UserUsecaseAdapter) GetByID(ctx context.Context, id string) (*authEntity.User, error) {
	user, err := a.userUsecase.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toAuthUser(user), nil
}

func (a *UserUsecaseAdapter) Create(ctx context.Context, user *authEntity.User) error {
	externalUser := toUserEntity(user)
	return a.userUsecase.Create(ctx, externalUser)
}

func (a *UserUsecaseAdapter) ComparePassword(user *authEntity.User, password string) bool {
	return a.userUsecase.ComparePassword(toUserEntity(user), password)
}

func (a *UserUsecaseAdapter) Update(ctx context.Context, user *authEntity.User) error {
	externalUser := toUserEntity(user)
	return a.userUsecase.Update(ctx, externalUser)
}

func toAuthUser(user *userEntity.User) *authEntity.User {
	if user == nil {
		return nil
	}
	return &authEntity.User{
		ID:          user.ID,
		Email:       user.Email,
		Password:    user.PasswordHash,
		PhoneNumber: user.PhoneNumber,
	}
}

func toUserEntity(user *authEntity.User) *userEntity.User {
	if user == nil {
		return nil
	}
	return &userEntity.User{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.Password,
		PhoneNumber:  user.PhoneNumber,
	}
}
