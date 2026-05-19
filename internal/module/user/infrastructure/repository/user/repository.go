package user_repository

import (
	"context"
	"time"

	user_entity "github.com/Fi44er/sdmed/internal/module/user/entity"
	user_model "github.com/Fi44er/sdmed/internal/module/user/infrastructure/repository/model"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type UserRepository struct {
	logger    *logger.Logger
	db        *gorm.DB
	converter *Converter
}

func NewUserRepository(logger *logger.Logger, db *gorm.DB) *UserRepository {
	return &UserRepository{
		logger:    logger,
		db:        db,
		converter: &Converter{},
	}
}

func (r *UserRepository) DeleteExpiredShadows(ctx context.Context) (int64, error) {
	r.logger.Info("Deleting expired shadow users")
	result := r.db.WithContext(ctx).
		Where("is_shadow = ? AND shadow_expires_at < ?", true, time.Now()).
		Delete(&user_model.User{})
	if result.Error != nil {
		r.logger.Errorf("Error deleting expired shadow users: %v", result.Error)
		return 0, result.Error
	}
	r.logger.Infof("Deleted %d expired shadow users", result.RowsAffected)
	return result.RowsAffected, nil
}

func (r *UserRepository) Create(ctx context.Context, user *user_entity.User) error {
	r.logger.Infof("Creating user: %+v", user)

	// Assign default role if no roles are explicitly specified
	if len(user.Roles) == 0 {
		var defaultRoleName string
		if user.IsShadow {
			defaultRoleName = "guest"
		} else {
			defaultRoleName = "user"
		}

		var dbRole user_model.Role
		if err := r.db.WithContext(ctx).Where("name = ?", defaultRoleName).First(&dbRole).Error; err == nil {
			user.Roles = append(user.Roles, user_entity.Role{
				ID:   dbRole.ID,
				Name: dbRole.Name,
			})
		} else {
			r.logger.Warnf("Default role %s not found in DB: %v", defaultRoleName, err)
		}
	}

	userModel := r.converter.ToModel(user)
	if err := r.db.WithContext(ctx).Create(userModel).Error; err != nil {
		r.logger.Errorf("Error creating user: %v", err)
		return err
	}
	user.ID = userModel.ID
	r.logger.Info("User created successfully")
	return nil
}

func (r *UserRepository) Promote(ctx context.Context, user *user_entity.User) error {
	// Promoted user is a real user. Assign the default "user" role.
	var dbRole user_model.Role
	if err := r.db.WithContext(ctx).Where("name = ?", "user").First(&dbRole).Error; err == nil {
		user.Roles = []user_entity.Role{
			{
				ID:   dbRole.ID,
				Name: dbRole.Name,
			},
		}
	} else {
		r.logger.Warnf("Default user role not found in DB: %v", err)
	}

	userModel := r.converter.ToModel(user)
	if err := r.db.WithContext(ctx).
		Model(&user_model.User{}).
		Where("id = ?", user.ID).
		Select(
			"email", "password_hash", "name", "surname",
			"patronymic", "phone_number",
			"is_shadow", "shadow_created_at", "shadow_expires_at",
			"updated_at",
		).
		Updates(userModel).Error; err != nil {
		r.logger.Errorf("Error promoting user: %v", err)
		return err
	}

	// Update the user_roles many-to-many association in the database
	if err := r.db.WithContext(ctx).Model(userModel).Association("Roles").Replace(userModel.Roles); err != nil {
		r.logger.Errorf("Error replacing user roles on promote: %v", err)
		return err
	}

	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *user_entity.User) error {
	r.logger.Infof("Updating user: %+v", user)
	userModel := r.converter.ToModel(user)
	if err := r.db.WithContext(ctx).Model(&user_model.User{}).Where("id = ?", user.ID).Updates(userModel).Error; err != nil {
		r.logger.Errorf("Error updating user: %v", err)
		return err
	}
	r.logger.Info("User updated successfully")
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	r.logger.Infof("Deleting user: %s", id)
	if err := r.db.WithContext(ctx).Delete(&user_model.User{}, "id = ?", id).Error; err != nil {
		r.logger.Errorf("Error deleting user: %v", err)
		return err
	}
	r.logger.Info("User deleted successfully")
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*user_entity.User, error) {
	r.logger.Infof("Getting user: %s", id)
	var userModel user_model.User
	if err := r.db.WithContext(ctx).Preload("Roles").First(&userModel, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("User not found: %s", id)
			return nil, nil
		}
		r.logger.Errorf("Error getting user: %v", err)
		return nil, err
	}
	user := r.converter.ToEntity(&userModel)
	r.logger.Info("User got successfully")
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*user_entity.User, error) {
	r.logger.Infof("Getting user by email: %s", email)
	var userModel user_model.User
	if err := r.db.WithContext(ctx).Preload("Roles").First(&userModel, "email = ?", email).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			r.logger.Warnf("User not found: %s", email)
			return nil, nil
		}
		r.logger.Errorf("Error getting user: %v", err)
		return nil, err
	}
	user := r.converter.ToEntity(&userModel)
	r.logger.Info("User got successfully")
	return user, nil
}

func (r *UserRepository) GetAll(ctx context.Context, limit int, offset int) ([]user_entity.User, error) {
	r.logger.Infof("Getting all users")
	var userModels []user_model.User
	if limit == 0 {
		limit = -1
	}
	if offset == 0 {
		offset = -1
	}
	if err := r.db.WithContext(ctx).Preload("Roles").Limit(limit).Offset(offset).Find(&userModels).Error; err != nil {
		r.logger.Errorf("Error getting users: %v", err)
		return nil, err
	}
	users := make([]user_entity.User, len(userModels))
	for i, userModel := range userModels {
		users[i] = *r.converter.ToEntity(&userModel)
	}
	r.logger.Info("Users got successfully")
	return users, nil
}
