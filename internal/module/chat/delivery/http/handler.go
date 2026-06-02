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

// CreateTicket godoc
// @Summary Create a new support ticket
// @Description Creates a new customer support ticket with the first message
// @Tags tickets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTicketDTO true "Ticket creation data"
// @Router /chat/tickets [post]
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

// SendMessage godoc
// @Summary Send a message to a ticket
// @Description Send a new message to an existing ticket
// @Tags tickets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Ticket ID"
// @Param request body dto.SendMessageDTO true "Message data"
// @Router /chat/tickets/{id}/messages [post]
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

// GetTicketByID godoc
// @Summary Get ticket by ID
// @Description Get detailed information about a specific ticket
// @Tags tickets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Ticket ID"
// @Router /chat/tickets/{id} [get]
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

// GetTicketMessages godoc
// @Summary Get ticket messages
// @Description Get paginated messages for a specific ticket
// @Tags tickets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Ticket ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Router /chat/tickets/{id}/messages [get]
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

// CloseTicket godoc
// @Summary Close a ticket
// @Description Close an existing ticket
// @Tags tickets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Ticket ID"
// @Router /chat/tickets/{id}/close [put]
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

// AssignTicket godoc
// @Summary Assign ticket to operator
// @Description Assign a ticket to the current user (operator)
// @Tags tickets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Ticket ID"
// @Router /chat/tickets/{id}/assign [put]
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

// MarkMessageAsRead godoc
// @Summary Mark message as read
// @Description Mark a specific message as read by the current user
// @Tags messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Message ID"
// @Router /chat/messages/{id}/read [post]
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

// GetMyTickets godoc
// @Summary Get my tickets
// @Description Get all tickets created by the current user (customer view)
// @Tags tickets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Router /chat/tickets/my [get]
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

// GetAllTickets godoc
// @Summary Get all tickets (admin/manager)
// @Description Get all tickets with pagination (admin and manager only)
// @Tags tickets
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Router /chat/tickets [get]
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
