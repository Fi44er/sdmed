package user_session_repository

import (
	auth_entity "github.com/Fi44er/sdmed/internal/module/auth/entity"
	auth_models "github.com/Fi44er/sdmed/internal/module/auth/infrastucture/repository/models"
)

type Converter struct {
}

func (c *Converter) ToEntity(session *auth_models.UserSession) *auth_entity.UserSession {
	return &auth_entity.UserSession{
		ID:        session.ID,
		UserID:    session.UserID,
		UserAgent: session.UserAgent,
		LastIP:    session.LastIP,

		ExpiresAt:  session.ExpiresAt,
		CreatedAt:  session.CreatedAt,
		UpdatedAt:  session.UpdatedAt,
		LastUsedAt: session.LastUsedAt,
	}
}

func (c *Converter) ToModel(session *auth_entity.UserSession) *auth_models.UserSession {
	return &auth_models.UserSession{
		ID:        session.ID,
		UserID:    session.UserID,
		UserAgent: session.UserAgent,
		LastIP:    session.LastIP,

		ExpiresAt:  session.ExpiresAt,
		CreatedAt:  session.CreatedAt,
		UpdatedAt:  session.UpdatedAt,
		LastUsedAt: session.LastUsedAt,
	}
}
