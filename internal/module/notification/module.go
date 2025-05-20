package notification

import (
	"github.com/Fi44er/sdmedik/backend/internal/config"
	"github.com/Fi44er/sdmedik/backend/internal/module/notification/service"
	"github.com/Fi44er/sdmedik/backend/internal/module/notification/service/smtp"
	"github.com/Fi44er/sdmedik/backend/pkg/logger"
)

type NotificationModule struct {
	logger *logger.Logger
	config *config.Config

	service *service.NotificationService
}

func NewNotificationModule(logger *logger.Logger, config *config.Config) *NotificationModule {
	return &NotificationModule{
		logger: logger,
		config: config,
	}
}

func (m *NotificationModule) Init() error {
	smtp, err := smtp.NewSMTPNotifier(
		m.config.SMTPHost,
		m.config.SMTPPort,
		m.config.SMTPFrom,
		m.config.SMTPPassword,
		m.config.SMTPTemplatePath+"index.html",
		5,
	)

	if err != nil {
		return err
	}

	notifiers := map[string]service.Notifier{
		"smtp": smtp,
	}

	m.service = service.NewNotificationService(notifiers)

	return nil
}
