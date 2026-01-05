package chat_ws

import (
	"encoding/json"

	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/gofiber/contrib/socketio"
)

type ChatWsHandler struct {
	logger *logger.Logger
}

func NewChatWsHandler(logger *logger.Logger) *ChatWsHandler {
	return &ChatWsHandler{
		logger: logger,
	}
}

type MsgObj struct {
	Event string `json:"event"`
	Data  any    `json:"data"`
}

func (h *ChatWsHandler) socketErr(ep *socketio.EventPayload, err error) {
	errByte := []byte(err.Error())
	ep.Kws.Emit(errByte)
}

func (h *ChatWsHandler) SetupSocketEvents() {
	socketio.On("connection", func(ep *socketio.EventPayload) {
		// USER_ONLINE отправляем что пользователь онлайн
		h.logger.Info("New connection established")
	})

	socketio.On("disconnect", func(ep *socketio.EventPayload) {
		// USER_OFFLINE отправляем что пользователь вышел из чата
		h.logger.Info("Connection closed")
	})

	socketio.On("chat_join", func(ep *socketio.EventPayload) {
		// USER_OFFLINE отправляем что пользователь вышел из чата
		h.logger.Info("Connection closed")
	})

	socketio.On("chat_leave", func(ep *socketio.EventPayload) {
		// USER_OFFLINE отправляем что пользователь вышел из чата
		h.logger.Info("Connection closed")
	})

	socketio.On("message_send", func(ep *socketio.EventPayload) {
		// server отправляет MESSAGE_NEW всем участникам чата
		h.logger.Info("Message received")
	})

	socketio.On("message_read", func(ep *socketio.EventPayload) {
		// server отправляет MESSAGE_READ_UPDATE отправителю сообщения для отображения прочтения
		h.logger.Info("Message received")
	})

	socketio.On("message_edit", func(ep *socketio.EventPayload) {
		h.logger.Info("Message received")
	})

	socketio.On("message_delete", func(ep *socketio.EventPayload) {
		h.logger.Info("Message received")
	})

	socketio.On("typing", func(ep *socketio.EventPayload) {
		// server отправляет USER_TYPING всем участникам чата
		h.logger.Info("Message received")
	})
}

func encodeMsg[T any](data any) (*T, error) {
	if data == nil {
		return nil, nil
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var res *T

	err = json.Unmarshal(dataBytes, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
