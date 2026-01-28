package user_session_repository

import (
	"context"
	"time"

	auth_entity "github.com/Fi44er/sdmed/internal/module/auth/entity"
	auth_models "github.com/Fi44er/sdmed/internal/module/auth/infrastucture/repository/models"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type IUserSessionRepository interface {
	Create(ctx context.Context, userSession *auth_entity.UserSession) error
	Delete(ctx context.Context, userSessionID string) error
	Get(ctx context.Context, userSessionID string) (*auth_entity.UserSession, error)
	Update(ctx context.Context, userSession *auth_entity.UserSession) error

	// NEW: Device management methods
	GetByUserID(ctx context.Context, userID string) ([]*auth_entity.UserSession, error)
	RevokeSession(ctx context.Context, sessionID string) error
	RevokeAllExcept(ctx context.Context, userID string, exceptSessionID string) error
	RevokeAll(ctx context.Context, userID string) error
	DeleteExpired(ctx context.Context) error
	IsRevoked(ctx context.Context, sessionID string) (bool, error)
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
	r.logger.Infof("Creating user session %s", userSession.ID)
	model := r.converter.ToModel(userSession)
	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		r.logger.Errorf("Failed to create user session: %v", err)
		return err
	}
	r.logger.Infof("User session %s created successfully", userSession.ID)
	return nil
}

func (r *UserSessionRepository) Delete(ctx context.Context, userSessionID string) error {
	r.logger.Infof("Deleting user session %s", userSessionID)
	if err := r.db.WithContext(ctx).Where("id = ?", userSessionID).Delete(&auth_models.UserSession{}).Error; err != nil {
		r.logger.Errorf("Failed to delete user session: %v", err)
		return err
	}
	r.logger.Infof("User session %s deleted successfully", userSessionID)
	return nil
}

func (r *UserSessionRepository) Update(ctx context.Context, userSession *auth_entity.UserSession) error {
	r.logger.Infof("Updating user session %s", userSession.ID)
	model := r.converter.ToModel(userSession)
	if err := r.db.WithContext(ctx).Model(&auth_models.UserSession{}).
		Where("id = ?", userSession.ID).
		Updates(model).Error; err != nil {
		r.logger.Errorf("Failed to update user session: %v", err)
		return err
	}
	r.logger.Infof("User session %s updated successfully", userSession.ID)
	return nil
}

func (r *UserSessionRepository) Get(ctx context.Context, userSessionID string) (*auth_entity.UserSession, error) {
	r.logger.Infof("Getting user session %s", userSessionID)
	var model auth_models.UserSession
	if err := r.db.WithContext(ctx).Where("id = ?", userSessionID).First(&model).Error; err != nil {
		r.logger.Errorf("Failed to get user session: %v", err)
		return nil, err
	}
	entity := r.converter.ToEntity(&model)
	r.logger.Infof("User session %s retrieved successfully", userSessionID)
	return entity, nil
}

// NEW: Get all sessions for a user
func (r *UserSessionRepository) GetByUserID(ctx context.Context, userID string) ([]*auth_entity.UserSession, error) {
	r.logger.Infof("Getting all sessions for user %s", userID)
	var models []auth_models.UserSession
	if err := r.db.WithContext(ctx).
		Where("user_id = ? AND is_revoked = ?", userID, false).
		Order("updated_at DESC").
		Find(&models).Error; err != nil {
		r.logger.Errorf("Failed to get user sessions: %v", err)
		return nil, err
	}

	sessions := make([]*auth_entity.UserSession, len(models))
	for i, model := range models {
		sessions[i] = r.converter.ToEntity(&model)
	}

	r.logger.Infof("Retrieved %d sessions for user %s", len(sessions), userID)
	return sessions, nil
}

// NEW: Revoke (ban) a specific session
func (r *UserSessionRepository) RevokeSession(ctx context.Context, sessionID string) error {
	r.logger.Infof("Revoking session %s", sessionID)
	if err := r.db.WithContext(ctx).
		Model(&auth_models.UserSession{}).
		Where("id = ?", sessionID).
		Update("is_revoked", true).Error; err != nil {
		r.logger.Errorf("Failed to revoke session: %v", err)
		return err
	}
	r.logger.Infof("Session %s revoked successfully", sessionID)
	return nil
}

// NEW: Revoke all sessions except current one
func (r *UserSessionRepository) RevokeAllExcept(ctx context.Context, userID string, exceptSessionID string) error {
	r.logger.Infof("Revoking all sessions for user %s except %s", userID, exceptSessionID)
	if err := r.db.WithContext(ctx).
		Model(&auth_models.UserSession{}).
		Where("user_id = ? AND id != ?", userID, exceptSessionID).
		Update("is_revoked", true).Error; err != nil {
		r.logger.Errorf("Failed to revoke sessions: %v", err)
		return err
	}
	r.logger.Infof("All sessions revoked for user %s except current", userID)
	return nil
}

// NEW: Revoke all sessions for a user
func (r *UserSessionRepository) RevokeAll(ctx context.Context, userID string) error {
	r.logger.Infof("Revoking all sessions for user %s", userID)
	if err := r.db.WithContext(ctx).
		Model(&auth_models.UserSession{}).
		Where("user_id = ?", userID).
		Update("is_revoked", true).Error; err != nil {
		r.logger.Errorf("Failed to revoke all sessions: %v", err)
		return err
	}
	r.logger.Infof("All sessions revoked for user %s", userID)
	return nil
}

// NEW: Delete expired sessions (for cleanup job)
func (r *UserSessionRepository) DeleteExpired(ctx context.Context) error {
	r.logger.Info("Deleting expired sessions")
	result := r.db.WithContext(ctx).
		Where("expires_at < ?", time.Now()).
		Delete(&auth_models.UserSession{})

	if result.Error != nil {
		r.logger.Errorf("Failed to delete expired sessions: %v", result.Error)
		return result.Error
	}

	r.logger.Infof("Deleted %d expired sessions", result.RowsAffected)
	return nil
}

// NEW: Check if session is revoked
func (r *UserSessionRepository) IsRevoked(ctx context.Context, sessionID string) (bool, error) {
	var session auth_models.UserSession
	if err := r.db.WithContext(ctx).
		Select("is_revoked").
		Where("id = ?", sessionID).
		First(&session).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return true, nil // session не найдена = считаем revoked
		}
		return false, err
	}
	return session.IsRevoked, nil
}
