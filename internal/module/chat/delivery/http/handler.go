package chat_http

import (
	"context"
	"strconv"

	"github.com/Fi44er/sdmed/internal/module/chat/dto"
	chat_entity "github.com/Fi44er/sdmed/internal/module/chat/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/session"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type IChatUsecase interface {
	CreateTicket(ctx context.Context, clientID string, subject string, priority int, messageText string) (*chat_entity.Chat, error)
	SendMessage(ctx context.Context, senderID string, isAdmin bool, isManager bool, chatID string, text string) (*chat_entity.Message, error)
	GetTicketByID(ctx context.Context, userID string, isAdmin bool, isManager bool, chatID string) (*chat_entity.Chat, error)
	GetTicketMessages(ctx context.Context, userID string, isAdmin bool, isManager bool, chatID string, page, pageSize int) ([]chat_entity.Message, error)
	CloseTicket(ctx context.Context, userID string, isAdmin bool, isManager bool, chatID string) error
	AssignTicket(ctx context.Context, operatorID string, chatID string) error
	MarkMessageAsRead(ctx context.Context, userID string, isAdmin bool, isManager bool, messageID string) error
	GetMyTickets(ctx context.Context, clientID string, page, pageSize int) ([]chat_entity.Chat, int64, error)
	GetOperatorOrAllTickets(ctx context.Context, userID string, isAdmin bool, isManager bool, page, pageSize int) ([]chat_entity.Chat, int64, error)
}

type ChatHandler struct {
	usecase   IChatUsecase
	logger    *logger.Logger
	validator *validator.Validate
}

func NewChatHandler(usecase IChatUsecase, logger *logger.Logger, validator *validator.Validate) *ChatHandler {
	return &ChatHandler{
		usecase:   usecase,
		logger:    logger,
		validator: validator,
	}
}

func (h *ChatHandler) isAdmin(c *fiber.Ctx) bool {
	sess := session.FromFiberContext(c)
	if sess == nil {
		return false
	}
	sessionData, ok := sess.Get("session_info").(map[string]any)
	if !ok || sessionData == nil {
		return false
	}

	if roles, ok := sessionData["user_roles"].([]any); ok {
		for _, r := range roles {
			if rStr, ok := r.(string); ok && rStr == "admin" {
				return true
			}
		}
	}
	return false
}

func (h *ChatHandler) isManager(c *fiber.Ctx) bool {
	sess := session.FromFiberContext(c)
	if sess == nil {
		return false
	}
	sessionData, ok := sess.Get("session_info").(map[string]any)
	if !ok || sessionData == nil {
		return false
	}

	if roles, ok := sessionData["user_roles"].([]any); ok {
		for _, r := range roles {
			if rStr, ok := r.(string); ok && rStr == "manager" {
				return true
			}
		}
	}
	return false
}

func (h *ChatHandler) getUserID(c *fiber.Ctx) string {
	userIDVal := c.Locals("user_id")
	if userID, ok := userIDVal.(string); ok {
		return userID
	}
	return ""
}

func (h *ChatHandler) CreateTicket(c *fiber.Ctx) error {
	h.logger.Info("HTTP: Creating ticket")
	userID := h.getUserID(c)

	var reqBody dto.CreateTicketDTO
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "invalid request body"})
	}

	if err := h.validator.Struct(&reqBody); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	chat, err := h.usecase.CreateTicket(c.Context(), userID, reqBody.Subject, reqBody.Priority, reqBody.FirstMessage)
	if err != nil {
		return err
	}

	resp := dto.TicketResponse{
		ID:         chat.ID,
		ClientID:   chat.ClientID,
		OperatorID: chat.OperatorID,
		Subject:    chat.Subject,
		Status:     string(chat.Status),
		Priority:   chat.Priority,
		Tags:       chat.Tags,
		CreatedAt:  chat.CreatedAt,
		UpdatedAt:  chat.UpdatedAt,
		ClosedAt:   chat.ClosedAt,
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "data": resp})
}

func (h *ChatHandler) SendMessage(c *fiber.Ctx) error {
	h.logger.Info("HTTP: Sending message")
	userID := h.getUserID(c)
	adminVal := h.isAdmin(c)
	mgrVal := h.isManager(c)
	chatID := c.Params("id")

	var reqBody dto.SendMessageDTO
	if err := c.BodyParser(&reqBody); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": "invalid request body"})
	}

	if err := h.validator.Struct(&reqBody); err != nil {
		return c.Status(400).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	msg, err := h.usecase.SendMessage(c.Context(), userID, adminVal, mgrVal, chatID, reqBody.Payload)
	if err != nil {
		return err
	}

	var replyToID *string
	if msg.ReplyToID != "" {
		replyToID = &msg.ReplyToID
	}

	resp := dto.MessageResponse{
		ID:        msg.ID,
		ChatID:    msg.ChatID,
		SenderID:  msg.SenderID,
		Type:      string(msg.Type),
		Payload:   msg.Payload,
		ReplyToID: replyToID,
		IsEdited:  msg.IsEdited,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
		ReadAt:    msg.ReadAt,
	}

	return c.Status(201).JSON(fiber.Map{"status": "success", "data": resp})
}

func (h *ChatHandler) GetTicketByID(c *fiber.Ctx) error {
	h.logger.Info("HTTP: Get ticket by ID")
	userID := h.getUserID(c)
	adminVal := h.isAdmin(c)
	mgrVal := h.isManager(c)
	chatID := c.Params("id")

	chat, err := h.usecase.GetTicketByID(c.Context(), userID, adminVal, mgrVal, chatID)
	if err != nil {
		return err
	}

	resp := dto.TicketResponse{
		ID:         chat.ID,
		ClientID:   chat.ClientID,
		OperatorID: chat.OperatorID,
		Subject:    chat.Subject,
		Status:     string(chat.Status),
		Priority:   chat.Priority,
		Tags:       chat.Tags,
		CreatedAt:  chat.CreatedAt,
		UpdatedAt:  chat.UpdatedAt,
		ClosedAt:   chat.ClosedAt,
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "data": resp})
}

func (h *ChatHandler) GetTicketMessages(c *fiber.Ctx) error {
	h.logger.Info("HTTP: Get messages")
	userID := h.getUserID(c)
	adminVal := h.isAdmin(c)
	mgrVal := h.isManager(c)
	chatID := c.Params("id")

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	msgs, err := h.usecase.GetTicketMessages(c.Context(), userID, adminVal, mgrVal, chatID, page, pageSize)
	if err != nil {
		return err
	}

	respList := make([]dto.MessageResponse, len(msgs))
	for i, msg := range msgs {
		var replyToID *string
		if msg.ReplyToID != "" {
			replyToID = &msg.ReplyToID
		}
		respList[i] = dto.MessageResponse{
			ID:        msg.ID,
			ChatID:    msg.ChatID,
			SenderID:  msg.SenderID,
			Type:      string(msg.Type),
			Payload:   msg.Payload,
			ReplyToID: replyToID,
			IsEdited:  msg.IsEdited,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
			ReadAt:    msg.ReadAt,
		}
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "data": respList})
}

func (h *ChatHandler) CloseTicket(c *fiber.Ctx) error {
	h.logger.Info("HTTP: Closing ticket")
	userID := h.getUserID(c)
	adminVal := h.isAdmin(c)
	mgrVal := h.isManager(c)
	chatID := c.Params("id")

	err := h.usecase.CloseTicket(c.Context(), userID, adminVal, mgrVal, chatID)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "ticket closed successfully"})
}

func (h *ChatHandler) AssignTicket(c *fiber.Ctx) error {
	h.logger.Info("HTTP: Assign ticket to operator")
	operatorID := h.getUserID(c)
	chatID := c.Params("id")

	err := h.usecase.AssignTicket(c.Context(), operatorID, chatID)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "ticket assigned successfully"})
}

func (h *ChatHandler) MarkMessageAsRead(c *fiber.Ctx) error {
	h.logger.Info("HTTP: Mark message as read")
	userID := h.getUserID(c)
	adminVal := h.isAdmin(c)
	mgrVal := h.isManager(c)
	messageID := c.Params("id")

	err := h.usecase.MarkMessageAsRead(c.Context(), userID, adminVal, mgrVal, messageID)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"status": "success", "message": "message marked as read successfully"})
}

func (h *ChatHandler) GetMyTickets(c *fiber.Ctx) error {
	h.logger.Info("HTTP: Get my tickets")
	userID := h.getUserID(c)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	chats, total, err := h.usecase.GetMyTickets(c.Context(), userID, page, pageSize)
	if err != nil {
		return err
	}

	respList := make([]dto.TicketResponse, len(chats))
	for i, chat := range chats {
		respList[i] = dto.TicketResponse{
			ID:         chat.ID,
			ClientID:   chat.ClientID,
			OperatorID: chat.OperatorID,
			Subject:    chat.Subject,
			Status:     string(chat.Status),
			Priority:   chat.Priority,
			Tags:       chat.Tags,
			CreatedAt:  chat.CreatedAt,
			UpdatedAt:  chat.UpdatedAt,
			ClosedAt:   chat.ClosedAt,
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   respList,
		"meta":   fiber.Map{"total": total, "page": page, "page_size": pageSize},
	})
}

func (h *ChatHandler) GetAllTickets(c *fiber.Ctx) error {
	h.logger.Info("HTTP: Get all tickets")
	userID := h.getUserID(c)
	adminVal := h.isAdmin(c)
	mgrVal := h.isManager(c)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	chats, total, err := h.usecase.GetOperatorOrAllTickets(c.Context(), userID, adminVal, mgrVal, page, pageSize)
	if err != nil {
		return err
	}

	respList := make([]dto.TicketResponse, len(chats))
	for i, chat := range chats {
		respList[i] = dto.TicketResponse{
			ID:         chat.ID,
			ClientID:   chat.ClientID,
			OperatorID: chat.OperatorID,
			Subject:    chat.Subject,
			Status:     string(chat.Status),
			Priority:   chat.Priority,
			Tags:       chat.Tags,
			CreatedAt:  chat.CreatedAt,
			UpdatedAt:  chat.UpdatedAt,
			ClosedAt:   chat.ClosedAt,
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   respList,
		"meta":   fiber.Map{"total": total, "page": page, "page_size": pageSize},
	})
}
