package chat_module

import (
	chat_http "github.com/Fi44er/sdmed/internal/module/chat/delivery/http"
	chat_ws "github.com/Fi44er/sdmed/internal/module/chat/delivery/ws"
	chat_repository "github.com/Fi44er/sdmed/internal/module/chat/infrastructure/repository/chat"
	message_repository "github.com/Fi44er/sdmed/internal/module/chat/infrastructure/repository/message"
	chat_usecase "github.com/Fi44er/sdmed/internal/module/chat/usecase"
	notification_service "github.com/Fi44er/sdmed/internal/module/notification/service"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type ChatModule struct {
	logger              *logger.Logger
	db                  *gorm.DB
	validator           *validator.Validate
	redisClient         *redis.Client
	userUsecase         chat_usecase.IUserUsecase
	notificationService *notification_service.NotificationService

	chatRepository    chat_repository.IChatRepository
	messageRepository message_repository.IMessageRepository
	chatUsecase       *chat_usecase.ChatUsecase
	chatHandler       *chat_http.ChatHandler
	chatWsHandler     *chat_ws.ChatWsHandler
}

func NewChatModule(
	logger *logger.Logger,
	db *gorm.DB,
	validator *validator.Validate,
	redisClient *redis.Client,
	userUsecase chat_usecase.IUserUsecase,
	notificationService *notification_service.NotificationService,
) *ChatModule {
	return &ChatModule{
		logger:              logger,
		db:                  db,
		validator:           validator,
		redisClient:         redisClient,
		userUsecase:         userUsecase,
		notificationService: notificationService,
	}
}

func (m *ChatModule) Init() {
	m.chatRepository = chat_repository.NewChatRepository(m.logger, m.db)
	m.messageRepository = message_repository.NewMessageRepository(m.logger, m.db)
	m.chatUsecase = chat_usecase.NewChatUsecase(
		m.logger,
		m.redisClient,
		m.chatRepository,
		m.messageRepository,
		m.userUsecase,
		m.notificationService,
	)
	m.chatHandler = chat_http.NewChatHandler(m.chatUsecase, m.logger, m.validator)
	m.chatWsHandler = chat_ws.NewChatWsHandler(m.logger, m.chatUsecase)
	m.chatWsHandler.SetupSocketEvents()
}

func (m *ChatModule) InitDelivery(router fiber.Router) {
	m.chatHandler.RegisterRoutes(router)
}
