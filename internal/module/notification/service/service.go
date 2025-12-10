package service

import (
	"time"

	"github.com/Fi44er/sdmed/pkg/logger"
)

type Message struct {
	Recipient    string
	Subject      string
	Content      string
	TemplatePath string
	Data         any
	Timestamp    time.Time
}

type Notifier interface {
	Send(msg *Message)
}

type NotificationService struct {
	notifiers map[string]Notifier
	logger    *logger.Logger
}

func NewNotificationService(notifiers map[string]Notifier, logger *logger.Logger) *NotificationService {
	return &NotificationService{
		notifiers: notifiers,
		logger:    logger,
	}
}

func (ns *NotificationService) Send(msg *Message, selectedNotifiers ...string) {
	for _, notifier := range selectedNotifiers {
		if notifier, ok := ns.notifiers[notifier]; ok {
			notifier.Send(msg)
		}
	}
}
