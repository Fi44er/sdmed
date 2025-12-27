package auth_adapters

import (
	"context"
	"strings"

	authEntity "github.com/Fi44er/sdmed/internal/module/auth/entity"
	constantAuth "github.com/Fi44er/sdmed/internal/module/auth/pkg/constant"
	userEntity "github.com/Fi44er/sdmed/internal/module/user/entity"
	constantUser "github.com/Fi44er/sdmed/internal/module/user/pkg/constant"
	role_usecase "github.com/Fi44er/sdmed/internal/module/user/usecase/role"
	userUsecase "github.com/Fi44er/sdmed/internal/module/user/usecase/user"
)

type UserUsecaseAdapter struct {
	userUsecase *userUsecase.UserUsecase
	roleUsecase role_usecase.IRoleUsecase
}

func NewUserUsecaseAdapter(userUsecase *userUsecase.UserUsecase, roleUsecase role_usecase.IRoleUsecase) *UserUsecaseAdapter {
	return &UserUsecaseAdapter{
		userUsecase: userUsecase,
		roleUsecase: roleUsecase,
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

func (a *UserUsecaseAdapter) GetAllRoles(ctx context.Context) ([]authEntity.Role, error) {
	roles, err := a.roleUsecase.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	rolesEntities := make([]authEntity.Role, len(roles))
	for i, role := range roles {
		rolesEntities[i] = *toAuthRole(&role)
	}
	return rolesEntities, nil
}

func toAuthUser(user *userEntity.User) *authEntity.User {
	if user == nil {
		return nil
	}

	roles := make([]authEntity.Role, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = *toAuthRole(&role)
	}

	return &authEntity.User{
		ID:          user.ID,
		Email:       user.Email,
		Password:    user.PasswordHash,
		PhoneNumber: user.PhoneNumber,
		Roles:       roles,
	}
}

func toAuthRole(role *userEntity.Role) *authEntity.Role {
	if role == nil {
		return nil
	}
	return &authEntity.Role{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: toAuthPermissions(role.Permissions),
	}
}

func toAuthPermissions(permissions []userEntity.Permission) []authEntity.Permission {
	authPermissions := make([]authEntity.Permission, len(permissions))
	for i, permission := range permissions {
		authPermissions[i] = toAuthPermission(permission)
	}
	return authPermissions
}

func toAuthPermission(permission userEntity.Permission) authEntity.Permission {
	return authEntity.Permission{
		ID:   permission.ID,
		Name: permission.Name,
	}
}

func toUserEntity(user *authEntity.User) *userEntity.User {
	if user == nil {
		return nil
	}

	name, surname, patronymic := splitFIO(user.FIO)
	return &userEntity.User{
		ID:           user.ID,
		Email:        user.Email,
		Name:         name,
		Surname:      surname,
		Patronymic:   patronymic,
		PasswordHash: user.Password,
		PhoneNumber:  user.PhoneNumber,
	}
}

func splitFIO(fio string) (name, surname, patronymic string) {
	parts := strings.Fields(fio)

	switch len(parts) {
	case 0:
		return "", "", ""
	case 1:
		return parts[0], "", ""
	case 2:
		return parts[0], parts[1], ""
	default:
		return parts[0], parts[1], strings.Join(parts[2:], " ")
	}
}
