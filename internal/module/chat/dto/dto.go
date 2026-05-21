package dto

import "time"

type CreateTicketDTO struct {
	Subject      string `json:"subject" validate:"required"`
	Priority     int    `json:"priority" validate:"required,min=1,max=3"`
	FirstMessage string `json:"first_message" validate:"required"`
}

type SendMessageDTO struct {
	Payload string `json:"payload" validate:"required"`
}

type TicketResponse struct {
	ID         string     `json:"id"`
	ClientID   string     `json:"client_id"`
	OperatorID *string    `json:"operator_id"`
	Subject    string     `json:"subject"`
	Status     string     `json:"status"`
	Priority   int        `json:"priority"`
	Tags       []string   `json:"tags"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	ClosedAt   *time.Time `json:"closed_at"`
}

type MessageResponse struct {
	ID        string     `json:"id"`
	ChatID    string     `json:"chat_id"`
	SenderID  string     `json:"sender_id"`
	Type      string     `json:"type"`
	Payload   string     `json:"payload"`
	ReplyToID *string    `json:"reply_to_id"`
	IsEdited  bool       `json:"is_edited"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	ReadAt    *time.Time `json:"read_at"`
}
