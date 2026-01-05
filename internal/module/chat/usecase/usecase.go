package chat_usecase

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/gofiber/contrib/socketio"
	"github.com/redis/go-redis/v9"
)

const (
	userSocketsPrefix      = "user_sockets:"
	roomParticipantsPrefix = "room_participants:"
	socketToUser           = "socket_to_user:"
	chatQueue              = "chat_queue"
)

type ICache interface {
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Get(ctx context.Context, key string, dest any) error
	Del(ctx context.Context, key string) error
}

type ChatUsecase struct {
	logger *logger.Logger
	cache  *redis.Client

	localSockets sync.Map
}

func NewChatUsecase(logger *logger.Logger, cache *redis.Client) *ChatUsecase {
	return &ChatUsecase{
		logger: logger,
		cache:  cache,
	}
}

func (u *ChatUsecase) AddConnection(ctx context.Context, userID string, k *socketio.Websocket) error {
	u.logger.Info("Connection to chat service established")

	socketID := k.GetUUID()
	u.localSockets.Store(socketID, k)

	if err := u.cache.SAdd(ctx, fmt.Sprintf("%s%s", userSocketsPrefix, userID), socketID).Err(); err != nil {
		u.logger.Errorf("failed to add socket to user: %v", err)
		return err
	}

	if err := u.cache.HSet(ctx, socketToUser, socketID, userID).Err(); err != nil {
		u.logger.Errorf("failed to add socket to user: %v", err)
		return err
	}

	return nil
}

func (u *ChatUsecase) JoinRoom(ctx context.Context, userID string, chatID string) error {
	if err := u.cache.SAdd(ctx, fmt.Sprintf("%s%s", roomParticipantsPrefix, chatID), userID).Err(); err != nil {
		u.logger.Errorf("failed to add user to room: %v", err)
		return err
	}

	return nil
}
