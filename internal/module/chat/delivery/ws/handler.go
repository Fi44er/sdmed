package chat_ws

import (
	"context"
	"encoding/json"

	chat_usecase "github.com/Fi44er/sdmed/internal/module/chat/usecase"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
)

type ChatWsHandler struct {
	logger  *logger.Logger
	usecase *chat_usecase.ChatUsecase
}

func NewChatWsHandler(logger *logger.Logger, usecase *chat_usecase.ChatUsecase) *ChatWsHandler {
	return &ChatWsHandler{
		logger:  logger,
		usecase: usecase,
	}
}

type MessageObject struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

type JoinPayload struct {
	UserID string `json:"user_id"`
	ChatID string `json:"chat_id"`
}

type SendMsgPayload struct {
	SenderID               string `json:"sender_id"`
	IsSenderManagerOrAdmin bool   `json:"is_sender_manager_or_admin"`
	ChatID                 string `json:"chat_id"`
	Payload                string `json:"payload"`
}

type TypingPayload struct {
	ChatID string `json:"chat_id"`
	UserID string `json:"user_id"`
}

type ReadPayload struct {
	UserID                 string `json:"user_id"`
	IsSenderManagerOrAdmin bool   `json:"is_sender_manager_or_admin"`
	MessageID              string `json:"message_id"`
}

func (h *ChatWsHandler) socketErr(ep *socketio.EventPayload, err error) {
	h.logger.Errorf("Socket error: %v", err)
	errByte, _ := json.Marshal(fiber.Map{"event": "error", "error": err.Error()})
	ep.Kws.Emit(errByte)
}

func (h *ChatWsHandler) SetupSocketEvents() {
	// Connection established
	socketio.On("connection", func(ep *socketio.EventPayload) {
		h.logger.Info("WS: New connection established")
	})

	// Connection closed
	socketio.On("disconnect", func(ep *socketio.EventPayload) {
		h.logger.Info("WS: Connection closed")
		h.usecase.RemoveConnection(ep.Kws.GetUUID())
	})

	// Message dispatcher
	socketio.On("message", func(ep *socketio.EventPayload) {
		h.logger.Infof("WS: Raw message received: %s", string(ep.Data))
		message := new(MessageObject)
		if err := json.Unmarshal(ep.Data, message); err != nil {
			h.socketErr(ep, err)
			return
		}

		if message.Event != "" {
			ep.Kws.Fire(message.Event, ep.Data)
		}
	})

	// Custom Event: chat_join
	socketio.On("chat_join", func(ep *socketio.EventPayload) {
		h.logger.Info("WS: chat_join event received")
		payload, err := decodePayload[JoinPayload](ep.Data)
		if err != nil {
			h.socketErr(ep, err)
			return
		}
		if payload == nil {
			return
		}

		h.usecase.AddConnection(payload.UserID, ep.Kws)

		resp, _ := json.Marshal(fiber.Map{
			"event": "joined",
			"data":  fiber.Map{"chat_id": payload.ChatID, "user_id": payload.UserID},
		})
		ep.Kws.Emit(resp)
		h.logger.Infof("WS: User %s joined successfully", payload.UserID)
	})

	// Custom Event: message_send
	socketio.On("message_send", func(ep *socketio.EventPayload) {
		h.logger.Info("WS: message_send event received")
		payload, err := decodePayload[SendMsgPayload](ep.Data)
		if err != nil {
			h.socketErr(ep, err)
			return
		}
		if payload == nil {
			return
		}

		ctx := context.Background()
		_, err = h.usecase.SendMessage(ctx, payload.SenderID, false, payload.IsSenderManagerOrAdmin, payload.ChatID, payload.Payload)
		if err != nil {
			h.socketErr(ep, err)
			return
		}
	})

	// Custom Event: message_read
	socketio.On("message_read", func(ep *socketio.EventPayload) {
		h.logger.Info("WS: message_read event received")
		payload, err := decodePayload[ReadPayload](ep.Data)
		if err != nil {
			h.socketErr(ep, err)
			return
		}
		if payload == nil {
			return
		}

		ctx := context.Background()
		err = h.usecase.MarkMessageAsRead(ctx, payload.UserID, false, payload.IsSenderManagerOrAdmin, payload.MessageID)
		if err != nil {
			h.socketErr(ep, err)
			return
		}
	})

	// Custom Event: typing
	socketio.On("typing", func(ep *socketio.EventPayload) {
		payload, err := decodePayload[TypingPayload](ep.Data)
		if err != nil {
			h.socketErr(ep, err)
			return
		}
		if payload == nil {
			return
		}

		go h.usecase.BroadcastToChat(payload.ChatID, "user_typing", fiber.Map{
			"chat_id": payload.ChatID,
			"user_id": payload.UserID,
		})
	})
}

type WSMessage[T any] struct {
	Event string `json:"event"`
	Data  T      `json:"data"`
}

func decodePayload[T any](data []byte) (*T, error) {
	var msg WSMessage[T]
	if err := json.Unmarshal(data, &msg); err != nil {
		return nil, err
	}
	return &msg.Data, nil
}
