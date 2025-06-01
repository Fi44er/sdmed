package module

import (
	"github.com/Fi44er/sdmed/internal/config"
	"github.com/Fi44er/sdmed/internal/module/notification/service"
	"github.com/Fi44er/sdmed/internal/module/notification/service/smtp"
	"github.com/Fi44er/sdmed/pkg/logger"
)

type NotificationModule struct {
	logger *logger.Logger
	config *config.Config

	Service *service.NotificationService
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
		5,
	)

	if err != nil {
		return err
	}

	notifiers := map[string]service.Notifier{
		"smtp": smtp,
	}

	m.Service = service.NewNotificationService(notifiers, m.logger)

	return nil
}
