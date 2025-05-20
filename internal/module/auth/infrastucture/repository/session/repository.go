package repository

import (
	"context"
	"fmt"

	"github.com/Fi44er/sdmedik/backend/internal/module/auth/entity"
	"github.com/Fi44er/sdmedik/backend/pkg/logger"
	"github.com/Fi44er/sdmedik/backend/pkg/session"
)

type SessionRepository struct {
	logger *logger.Logger
}

func NewSessionRepository(
	logger *logger.Logger,
) *SessionRepository {
	return &SessionRepository{logger: logger}
}

func (r *SessionRepository) Get(ctx context.Context, key string) (*entity.UserSesion, error) {
	session, ok := ctx.Value("session").(session.Session)
	if !ok {
		return nil, fmt.Errorf("!")
	}

	_ = session.Get(key)
	return nil, nil
}
