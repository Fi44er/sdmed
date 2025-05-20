package service

import (
	"fmt"
	"time"

	"github.com/Fi44er/sdmedik/backend/pkg/logger"
)

type Message struct {
	Recipient string
	Subject   string
	Content   string
	Data      interface{}
	Timestamp time.Time
}

type Notifier interface {
	Send(msg *Message) error
}

type NotificationService struct {
	notifiers map[string]Notifier
	logger    *logger.Logger
}

func NewNotificationService(notifiers map[string]Notifier) *NotificationService {
	return &NotificationService{
		notifiers: notifiers,
	}
}

func (ns *NotificationService) Send(msg *Message, selectedNotifiers ...string) error {
	var errors []error
	for _, notifier := range selectedNotifiers {
		if notifier, ok := ns.notifiers[notifier]; ok {
			if err := notifier.Send(msg); err != nil {
				errors = append(errors, err)
			}
		}
	}
	if len(errors) > 0 {
		ns.logger.Errorf("ошибки при отправке: %v", errors)
		return fmt.Errorf("ошибки при отправке: %v", errors)
	}
	return nil
}
