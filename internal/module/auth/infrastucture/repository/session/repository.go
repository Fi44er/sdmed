package repository

import (
	"context"
	"fmt"

	"github.com/Fi44er/sdmed/internal/module/auth/entity"
	"github.com/Fi44er/sdmed/internal/module/auth/pkg/constant"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/session"
	"github.com/go-viper/mapstructure/v2"
)

type SessionRepository struct {
	logger *logger.Logger
}

func NewSessionRepository(
	logger *logger.Logger,
) *SessionRepository {
	return &SessionRepository{logger: logger}
}

const (
	sessionInfoKey = "session_info"
)

func (r *SessionRepository) GetSessionInfo(ctx context.Context) (*entity.UserSession, error) {
	session, ok := ctx.Value("session").(session.Session)
	if !ok {
		r.logger.Error("session not found")
		return nil, fmt.Errorf("session not found")
	}

	sessionData, ok := session.Get(sessionInfoKey).(map[string]interface{})
	if !ok {
		return nil, constant.ErrSessionInfoNotFound
	}

	var userSession entity.UserSession
	if err := mapstructure.Decode(sessionData, &userSession); err != nil {
		return nil, fmt.Errorf("failed to decode session data: %v", err)
	}

	return &userSession, nil
}

func (r *SessionRepository) PutSessionInfo(ctx context.Context, sessionInfo *entity.UserSession) error {
	session, ok := ctx.Value("session").(session.Session)
	if !ok {
		r.logger.Error("session not found")
		return constant.ErrSessionInfoNotFound
	}

	session.Put(sessionInfoKey, sessionInfo)
	return nil
}

func (r *SessionRepository) DeleteSessionInfo(ctx context.Context) error {
	session, ok := ctx.Value("session").(session.Session)
	if !ok {
		r.logger.Error("session not found")
		return constant.ErrSessionInfoNotFound
	}

	session.Delete(sessionInfoKey)
	return nil
}
