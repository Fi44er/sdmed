package chat_http

import (
	"github.com/Fi44er/sdmed/internal/middlewares"
	"github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
)

func (h *ChatHandler) RegisterRoutes(router fiber.Router) {
	// WebSocket endpoint
	router.Get("/chat/ws", socketio.New(func(kws *socketio.Websocket) {
		// Events are handled automatically by SetupSocketEvents package-level listeners
	}))

	tickets := router.Group("/chat/tickets", middlewares.RequireAuth())

	// Customer ticket routes
	tickets.Post("/", middlewares.Authorize("tickets", "my"), h.CreateTicket)
	tickets.Get("/my", middlewares.Authorize("tickets", "my"), h.GetMyTickets)

	// Manager / Operator ticket routes
	tickets.Get("/", middlewares.Authorize("tickets", "all"), h.GetAllTickets)
	tickets.Put("/:id/assign", middlewares.Authorize("tickets", "all"), h.AssignTicket)

	// Shared ticket details, messages, sending & closing routes
	tickets.Get("/:id", h.GetTicketByID)
	tickets.Get("/:id/messages", h.GetTicketMessages)
	tickets.Post("/:id/messages", h.SendMessage)
	tickets.Put("/:id/close", h.CloseTicket)

	// Message level routes
	messages := router.Group("/chat/messages", middlewares.RequireAuth())
	messages.Post("/:id/read", h.MarkMessageAsRead)
}
