package contracts

import (
	"context"
	"time"

	auth_entity "github.com/Fi44er/sdmed/internal/module/auth/entity"
	"github.com/Fi44er/sdmed/internal/module/notification/service"
)

type IUserUsecase interface {
	GetByEmail(ctx context.Context, email string) (*auth_entity.User, error)
	GetByID(ctx context.Context, id string) (*auth_entity.User, error)
	Create(ctx context.Context, user *auth_entity.User) error
	ComparePassword(user *auth_entity.User, password string) bool
	Update(ctx context.Context, user *auth_entity.User) error
}

type ICache interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string, dest any) error
	Del(ctx context.Context, key string) error
}

type INotificationService interface {
	Send(msg *service.Message, selectedNotifiers ...string)
}

type ISessionRepository interface {
	GetSessionInfo(ctx context.Context) (*auth_entity.ActiveSession, error)
	PutSessionInfo(ctx context.Context, sessionInfo *auth_entity.ActiveSession) error
	DeleteSessionInfo(ctx context.Context) error
}

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

type IShadowUserService interface {
	CreateShadowUser(ctx context.Context) (*auth_entity.User, error)
	PromoteToRealUser(ctx context.Context, shadowUserID string, user *auth_entity.User) error
	CleanupExpiredShadows(ctx context.Context) error
}
