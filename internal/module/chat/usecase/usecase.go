package chat_usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	chat_entity "github.com/Fi44er/sdmed/internal/module/chat/entity"
	chat_repository "github.com/Fi44er/sdmed/internal/module/chat/infrastructure/repository/chat"
	message_repository "github.com/Fi44er/sdmed/internal/module/chat/infrastructure/repository/message"
	chat_constant "github.com/Fi44er/sdmed/internal/module/chat/pkg/constant"
	notification_service "github.com/Fi44er/sdmed/internal/module/notification/service"
	user_entity "github.com/Fi44er/sdmed/internal/module/user/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type IUserUsecase interface {
	GetByID(ctx context.Context, id string) (*user_entity.User, error)
	GetByRole(ctx context.Context, roleName string) ([]user_entity.User, error)
}

type ChatUsecase struct {
	logger              *logger.Logger
	cache               *redis.Client
	chatRepo            chat_repository.IChatRepository
	messageRepo         message_repository.IMessageRepository
	userUsecase         IUserUsecase
	notificationService *notification_service.NotificationService

	localSockets sync.Map // Map[userID]*sync.Map (Map[socketUUID]*socketio.Websocket)
}

func NewChatUsecase(
	logger *logger.Logger,
	cache *redis.Client,
	chatRepo chat_repository.IChatRepository,
	messageRepo message_repository.IMessageRepository,
	userUsecase IUserUsecase,
	notificationService *notification_service.NotificationService,
) *ChatUsecase {
	return &ChatUsecase{
		logger:              logger,
		cache:               cache,
		chatRepo:            chatRepo,
		messageRepo:         messageRepo,
		userUsecase:         userUsecase,
		notificationService: notificationService,
	}
}

// Connection Management

func (u *ChatUsecase) AddConnection(userID string, k *socketio.Websocket) {
	u.logger.Infof("WS: User %s connected with socket %s", userID, k.GetUUID())
	actual, _ := u.localSockets.LoadOrStore(userID, &sync.Map{})
	userSockets := actual.(*sync.Map)
	userSockets.Store(k.GetUUID(), k)
}

func (u *ChatUsecase) RemoveConnection(socketID string) {
	u.localSockets.Range(func(key, value any) bool {
		userID := key.(string)
		userSockets := value.(*sync.Map)
		if _, ok := userSockets.Load(socketID); ok {
			userSockets.Delete(socketID)
			u.logger.Infof("WS: Removed connection %s for user %s", socketID, userID)
			return false // stop range iteration
		}
		return true
	})
}

func (u *ChatUsecase) EmitToUser(userID string, event string, data any) {
	actual, ok := u.localSockets.Load(userID)
	if !ok {
		return
	}
	userSockets := actual.(*sync.Map)

	resp, err := json.Marshal(fiber.Map{
		"event": event,
		"data":  data,
	})
	if err != nil {
		u.logger.Errorf("Failed to marshal WS event: %v", err)
		return
	}

	userSockets.Range(func(key, value any) bool {
		kws := value.(*socketio.Websocket)
		kws.Emit(resp)
		return true
	})
}

func (u *ChatUsecase) EmitToManagers(event string, data any) {
	managers, err := u.userUsecase.GetByRole(context.Background(), "manager")
	if err != nil {
		return
	}
	admins, _ := u.userUsecase.GetByRole(context.Background(), "admin")

	allOperators := append(managers, admins...)
	for _, op := range allOperators {
		u.EmitToUser(op.ID, event, data)
	}
}

func (u *ChatUsecase) BroadcastToChat(chatID string, event string, data any) {
	chat, err := u.chatRepo.GetByID(context.Background(), chatID)
	if err != nil || chat == nil {
		return
	}

	// Always emit to the client/customer
	u.EmitToUser(chat.ClientID, event, data)

	// Emit to assigned operator
	if chat.OperatorID != nil {
		u.EmitToUser(*chat.OperatorID, event, data)
	} else {
		// If no operator assigned yet, broadcast to all managers/admins
		u.EmitToManagers(event, data)
	}
}

// Support Chat Business Logic

func (u *ChatUsecase) checkTicketAccess(chat *chat_entity.Chat, userID string, isAdmin bool, isManager bool) bool {
	if isAdmin {
		return true
	}
	if isManager {
		// Manager only has access if they are assigned OR if the ticket is unassigned
		return chat.OperatorID == nil || *chat.OperatorID == userID
	}
	// Customer only has access to their own tickets
	return chat.ClientID == userID
}

func (u *ChatUsecase) CreateTicket(ctx context.Context, clientID string, subject string, priority int, messageText string) (*chat_entity.Chat, error) {
	u.logger.Infof("Usecase: Creating ticket for client %s, subject: %s", clientID, subject)

	chat := &chat_entity.Chat{
		ClientID:  clientID,
		Subject:   subject,
		Status:    chat_entity.StatusNew,
		Priority:  priority,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Try auto-assigning a manager round-robin or first available
	managers, err := u.userUsecase.GetByRole(ctx, "manager")
	if err == nil && len(managers) > 0 {
		assignedManager := managers[0]
		chat.OperatorID = &assignedManager.ID
		chat.Status = chat_entity.StatusOpen
	}

	if err := u.chatRepo.Create(ctx, chat); err != nil {
		return nil, err
	}

	// Create initial message
	msg := &chat_entity.Message{
		ChatID:    chat.ID,
		SenderID:  clientID,
		Type:      chat_entity.MsgTypeRegular,
		Payload:   messageText,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.messageRepo.Create(ctx, msg); err != nil {
		return nil, err
	}

	// Notify assigned manager or all managers if none assigned
	go u.notifyManagers(context.Background(), clientID, chat.OperatorID, subject, messageText)

	// Emit WS updates
	if chat.OperatorID != nil {
		go u.EmitToUser(*chat.OperatorID, "ticket_created", chat)
	} else {
		go u.EmitToManagers("ticket_created", chat)
	}

	return chat, nil
}

func (u *ChatUsecase) SendMessage(ctx context.Context, senderID string, isAdmin bool, isManager bool, chatID string, text string) (*chat_entity.Message, error) {
	u.logger.Infof("Usecase: Sending message in chat %s by sender %s", chatID, senderID)

	chat, err := u.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	if chat == nil {
		return nil, chat_constant.ErrChatNotFound
	}

	// Check if chat is closed
	if chat.Status == chat_entity.StatusClosed || chat.Status == chat_entity.StatusArchived {
		return nil, chat_constant.ErrChatClosed
	}

	// Enforce strict access permissions
	if !u.checkTicketAccess(chat, senderID, isAdmin, isManager) {
		return nil, chat_constant.ErrUnauthorized
	}

	msg := &chat_entity.Message{
		ChatID:    chatID,
		SenderID:  senderID,
		Type:      chat_entity.MsgTypeRegular,
		Payload:   text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := u.messageRepo.Create(ctx, msg); err != nil {
		return nil, err
	}

	// If manager replies, auto-assign if not assigned
	if isManager || isAdmin {
		if chat.OperatorID == nil || *chat.OperatorID != senderID {
			chat.OperatorID = &senderID
			chat.Status = chat_entity.StatusOpen
		}
	}

	chat.UpdatedAt = time.Now()
	if err := u.chatRepo.Update(ctx, chat); err != nil {
		u.logger.Errorf("Failed to update chat update time: %v", err)
	}

	// If customer sent the message, notify the operator/managers
	if !isManager && !isAdmin {
		go u.notifyManagers(context.Background(), senderID, chat.OperatorID, chat.Subject, text)
	}

	// Emit WebSocket event to chat participants
	go u.BroadcastToChat(chatID, "message_new", msg)

	return msg, nil
}

func (u *ChatUsecase) GetTicketByID(ctx context.Context, userID string, isAdmin bool, isManager bool, chatID string) (*chat_entity.Chat, error) {
	chat, err := u.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	if chat == nil {
		return nil, chat_constant.ErrChatNotFound
	}

	if !u.checkTicketAccess(chat, userID, isAdmin, isManager) {
		return nil, chat_constant.ErrUnauthorized
	}

	return chat, nil
}

func (u *ChatUsecase) GetTicketMessages(ctx context.Context, userID string, isAdmin bool, isManager bool, chatID string, page, pageSize int) ([]chat_entity.Message, error) {
	chat, err := u.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	if chat == nil {
		return nil, chat_constant.ErrChatNotFound
	}

	if !u.checkTicketAccess(chat, userID, isAdmin, isManager) {
		return nil, chat_constant.ErrUnauthorized
	}

	return u.messageRepo.GetByChatID(ctx, chatID, page, pageSize)
}

func (u *ChatUsecase) CloseTicket(ctx context.Context, userID string, isAdmin bool, isManager bool, chatID string) error {
	chat, err := u.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return err
	}
	if chat == nil {
		return chat_constant.ErrChatNotFound
	}

	if !u.checkTicketAccess(chat, userID, isAdmin, isManager) {
		return chat_constant.ErrUnauthorized
	}

	now := time.Now()
	chat.Status = chat_entity.StatusClosed
	chat.ClosedAt = &now
	chat.UpdatedAt = now

	err = u.chatRepo.Update(ctx, chat)
	if err == nil {
		go u.BroadcastToChat(chatID, "ticket_closed", chat)
	}
	return err
}

func (u *ChatUsecase) AssignTicket(ctx context.Context, operatorID string, chatID string) error {
	chat, err := u.chatRepo.GetByID(ctx, chatID)
	if err != nil {
		return err
	}
	if chat == nil {
		return chat_constant.ErrChatNotFound
	}

	chat.OperatorID = &operatorID
	chat.Status = chat_entity.StatusOpen
	chat.UpdatedAt = time.Now()

	err = u.chatRepo.Update(ctx, chat)
	if err == nil {
		go u.BroadcastToChat(chatID, "ticket_assigned", chat)
	}
	return err
}

func (u *ChatUsecase) MarkMessageAsRead(ctx context.Context, userID string, isAdmin bool, isManager bool, messageID string) error {
	msg, err := u.messageRepo.GetByID(ctx, messageID)
	if err != nil {
		return err
	}
	if msg == nil {
		return chat_constant.ErrMessageNotFound
	}

	chat, err := u.chatRepo.GetByID(ctx, msg.ChatID)
	if err != nil {
		return err
	}
	if chat == nil {
		return chat_constant.ErrChatNotFound
	}

	if !u.checkTicketAccess(chat, userID, isAdmin, isManager) {
		return chat_constant.ErrUnauthorized
	}

	if msg.SenderID != userID {
		err = u.messageRepo.MarkAsRead(ctx, messageID)
		if err == nil {
			go u.BroadcastToChat(msg.ChatID, "message_read_update", fiber.Map{"message_id": messageID})
		}
		return err
	}

	return nil
}

func (u *ChatUsecase) GetMyTickets(ctx context.Context, clientID string, page, pageSize int) ([]chat_entity.Chat, int64, error) {
	return u.chatRepo.GetByClientID(ctx, clientID, page, pageSize)
}

func (u *ChatUsecase) GetOperatorOrAllTickets(ctx context.Context, userID string, isAdmin bool, isManager bool, page, pageSize int) ([]chat_entity.Chat, int64, error) {
	if isAdmin {
		return u.chatRepo.GetAll(ctx, page, pageSize)
	}
	if isManager {
		return u.chatRepo.GetByOperatorID(ctx, userID, page, pageSize)
	}
	return nil, 0, chat_constant.ErrUnauthorized
}

// Notification Helper

func (u *ChatUsecase) notifyManagers(ctx context.Context, clientID string, operatorID *string, subject string, content string) {
	u.logger.Info("Background: Preparing manager notification")

	clientUser, err := u.userUsecase.GetByID(ctx, clientID)
	var clientName string
	if err != nil || clientUser == nil {
		clientName = "Пользователь"
	} else {
		clientName = fmt.Sprintf("%s %s", clientUser.Name, clientUser.Surname)
	}

	var targets []user_entity.User
	if operatorID != nil {
		opUser, err := u.userUsecase.GetByID(ctx, *operatorID)
		if err == nil && opUser != nil {
			targets = append(targets, *opUser)
		}
	} else {
		managers, err := u.userUsecase.GetByRole(ctx, "manager")
		if err == nil {
			targets = managers
		}
	}

	for _, mgr := range targets {
		if mgr.Email != "" {
			msg := &notification_service.Message{
				Recipient: mgr.Email,
				Subject:   "Новое сообщение в поддержке",
				Content:   fmt.Sprintf("Отправитель: %s\nТема: %s\n\nСообщение:\n%s", clientName, subject, content),
			}
			u.notificationService.Send(msg, "system", "smtp")

			// Emits real-time notification to active connections
			u.EmitToUser(mgr.ID, "system_notification", msg)
		}
	}
}
