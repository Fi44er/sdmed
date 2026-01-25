package user_session_repository

import (
	"context"

	auth_entity "github.com/Fi44er/sdmed/internal/module/auth/entity"
	auth_models "github.com/Fi44er/sdmed/internal/module/auth/infrastucture/repository/models"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type IUserSessionRepository interface {
	Create(ctx context.Context, userSession *auth_entity.UserSession) error
	Update(ctx context.Context, userSession *auth_entity.UserSession) error
	Delete(ctx context.Context, userSession *auth_entity.UserSession) error
	FindByID(ctx context.Context, id string) (*auth_entity.UserSession, error)
	FindByUserID(ctx context.Context, userID string) ([]auth_entity.UserSession, error)
}

type UserSessionRepository struct {
	logger    *logger.Logger
	db        *gorm.DB
	converter *Converter
}

func NewUserSessionRepository(logger *logger.Logger, db *gorm.DB) IUserSessionRepository {
	return &UserSessionRepository{
		logger:    logger,
		db:        db,
		converter: &Converter{},
	}
}

func (r *UserSessionRepository) Create(ctx context.Context, userSession *auth_entity.UserSession) error {
	r.logger.Infof("Creating user session: %v", userSession)
	model := r.converter.ToModel(userSession)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		r.logger.Errorf("Failed to create user session: %v", err)
		return err
	}

	r.logger.Infof("User session created: %v", model)
	return nil
}

func (r *UserSessionRepository) Update(ctx context.Context, userSession *auth_entity.UserSession) error {
	r.logger.Infof("Updating user session: %v", userSession)
	model := r.converter.ToModel(userSession)
	if err := r.db.WithContext(ctx).Model(model).Updates(model).Error; err != nil {
		r.logger.Errorf("Failed to update user session: %v", err)
		return err
	}

	r.logger.Infof("User session updated: %v", model)
	return nil
}

func (r *UserSessionRepository) Delete(ctx context.Context, userSession *auth_entity.UserSession) error {
	r.logger.Infof("Deleting user session: %v", userSession)
	model := r.converter.ToModel(userSession)
	if err := r.db.WithContext(ctx).Delete(model).Error; err != nil {
		r.logger.Errorf("Failed to delete user session: %v", err)
		return err
	}

	r.logger.Infof("User session deleted: %v", model)
	return nil
}

func (r *UserSessionRepository) FindByID(ctx context.Context, id string) (*auth_entity.UserSession, error) {
	r.logger.Infof("Finding user session by ID: %s", id)
	var model auth_models.UserSession
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&model).Error; err != nil {
		r.logger.Errorf("Failed to find user session by ID: %v", err)
		return nil, err
	}

	r.logger.Infof("User session found: %v", model)
	return r.converter.ToEntity(&model), nil
}

func (r *UserSessionRepository) FindByUserID(ctx context.Context, userID string) ([]auth_entity.UserSession, error) {
	r.logger.Infof("Finding user sessions by user ID: %s", userID)
	var models []auth_models.UserSession
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error; err != nil {
		r.logger.Errorf("Failed to find user sessions by user ID: %v", err)
		return nil, err
	}

	r.logger.Infof("User sessions found: %v", models)
	entities := make([]auth_entity.UserSession, len(models))
	for i, model := range models {
		entities[i] = *r.converter.ToEntity(&model)
	}
	return entities, nil
}
