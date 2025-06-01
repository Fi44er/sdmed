package repository

import (
	"context"
	"fmt"

	"github.com/Fi44er/sdmed/internal/module/auth/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/session"
)

type SessionRepository struct {
	logger *logger.Logger
}

func NewSessionRepository(
	logger *logger.Logger,
) *SessionRepository {
	return &SessionRepository{logger: logger}
}

func (r *SessionRepository) GetSessionInfo(ctx context.Context) (*entity.UserSesion, error) {
	session, ok := ctx.Value("session").(session.Session)
	if !ok {
		return nil, fmt.Errorf("session not found")
	}

	sessionInfo, ok := session.Get("session_info").(*entity.UserSesion)
	if !ok {
		return nil, fmt.Errorf("session info not found")
	}

	return sessionInfo, nil
}

func (r *SessionRepository) PutSessionInfo(ctx context.Context, sessionInfo *entity.UserSesion) error {
	session, ok := ctx.Value("session").(session.Session)
	if !ok {
		return fmt.Errorf("session not found")
	}

	session.Put("session_info", sessionInfo)
	return nil
}

func (r *SessionRepository) DeleteSessionInfo(ctx context.Context) error {
	session, ok := ctx.Value("session").(session.Session)
	if !ok {
		return fmt.Errorf("session not found")
	}

	session.Delete("session_info")
	return nil
}
