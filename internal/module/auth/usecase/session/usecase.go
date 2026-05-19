package session_usecase

import (
	"context"

	auth_entity "github.com/Fi44er/sdmed/internal/module/auth/entity"
)

type ISessionRepository interface {
	GetSessionInfo(ctx context.Context) (*auth_entity.ActiveSession, error)
}

type SessionUseCase struct {
	sessionRepo ISessionRepository
}

func NewSessionUseCase(sessionRepo ISessionRepository) *SessionUseCase {
	return &SessionUseCase{
		sessionRepo: sessionRepo,
	}
}

func (uc *SessionUseCase) GetSessionInfo(ctx context.Context) (*auth_entity.ActiveSession, error) {
	return uc.sessionRepo.GetSessionInfo(ctx)
}
