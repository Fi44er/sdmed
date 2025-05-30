package service

import (
	"fmt"
	"time"

	"github.com/Fi44er/sdmedik/backend/pkg/logger"
)

type Message struct {
	Recipient    string
	Subject      string
	Content      string
	TemplatePath string
	Data         interface{}
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

func (ns *NotificationService) Send(msg *Message, selectedNotifiers ...string) error {
	var errors []error
	for _, notifier := range selectedNotifiers {
		if notifier, ok := ns.notifiers[notifier]; ok {
			notifier.Send(msg)
		}
	}
	if len(errors) > 0 {
		ns.logger.Errorf("ошибки при отправке: %v", errors)
		return fmt.Errorf("ошибки при отправке: %v", errors)
	}
	return nil
}
