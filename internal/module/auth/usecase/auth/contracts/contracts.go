package contracts

import (
	"context"
	"time"

	"github.com/Fi44er/sdmed/internal/module/auth/entity"
	"github.com/Fi44er/sdmed/internal/module/notification/service"
)

type IUserUsecase interface {
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByID(ctx context.Context, id string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	ComparePassword(user *entity.User, password string) bool
}

type ICache interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Del(ctx context.Context, key string) error
}

type INotificationService interface {
	Send(msg *service.Message, selectedNotifiers ...string)
}

type ISessionRepository interface {
	GetSessionInfo(ctx context.Context) (*entity.UserSession, error)
	PutSessionInfo(ctx context.Context, sessionInfo *entity.UserSession) error
	DeleteSessionInfo(ctx context.Context) error
}

type ITokenService interface {
	CreateToken(userID string, ttl time.Duration, privateKey string) (*entity.TokenDetails, error)
	ValidateToken(token, publicKey string) (*entity.TokenDetails, error)
}
