package chat_entity

import (
	"time"

	chat_constant "github.com/Fi44er/sdmed/internal/module/chat/pkg/constant"
)

type ParticipantRole string

const (
	ParticipantRoleCustomer   = "customer"
	ParticipantRoleAgent      = "agent"
	ParticipantRoleSupervisor = "supervisor"
)

// Если в одном чате может быть > 2 человек (например, два оператора)
type ChatParticipant struct {
	ChatID   string
	UserID   string
	Role     ParticipantRole // "customer", "agent", "supervisor"
	JoinedAt time.Time
}

func NewChatParticipant(chatID, userID string, role ParticipantRole) *ChatParticipant {
	return &ChatParticipant{
		ChatID:   chatID,
		UserID:   userID,
		Role:     role,
		JoinedAt: time.Now(),
	}
}

func (e *ChatParticipant) ValidateRole() error {
	switch e.Role {
	case ParticipantRoleCustomer, ParticipantRoleAgent, ParticipantRoleSupervisor:
		return nil
	default:
		return chat_constant.ErrInvalidRole
	}
}
