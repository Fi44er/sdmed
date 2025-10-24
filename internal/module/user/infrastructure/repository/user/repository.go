package user_repository

import (
	"context"

	"github.com/Fi44er/sdmed/internal/module/user/entity"
	"github.com/Fi44er/sdmed/internal/module/user/infrastructure/repository/model"
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

func (r *UserRepository) Create(ctx context.Context, user *user_entity.User) error {
	r.logger.Infof("Creating user: %+v", user)
	userModel := r.converter.ToModel(user)
	if err := r.db.WithContext(ctx).Create(userModel).Error; err != nil {
		r.logger.Errorf("Error creating user: %v", err)
		return err
	}
	user.ID = userModel.ID
	r.logger.Info("User created successfully")
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
	if err := r.db.WithContext(ctx).First(&userModel, "id = ?", id).Error; err != nil {
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
	if err := r.db.WithContext(ctx).First(&userModel, "email = ?", email).Error; err != nil {
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
	if err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&userModels).Error; err != nil {
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
