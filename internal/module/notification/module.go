package module

import (
	"github.com/Fi44er/sdmed/internal/config"
	"github.com/Fi44er/sdmed/internal/module/notification/service"
	"github.com/Fi44er/sdmed/internal/module/notification/service/smtp"
	"github.com/Fi44er/sdmed/internal/module/notification/service/system"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type NotificationModule struct {
	logger *logger.Logger
	config *config.Config
	db     *gorm.DB

	service *service.NotificationService
}

func NewNotificationModule(logger *logger.Logger, config *config.Config, db *gorm.DB) *NotificationModule {
	return &NotificationModule{
		logger: logger,
		config: config,
		db:     db,
	}
}

func (m *NotificationModule) Init() error {
	smtpNotifier, err := smtp.NewSMTPNotifier(
		m.config.SMTPHost,
		m.config.SMTPPort,
		m.config.SMTPFrom,
		m.config.SMTPPassword,
		5,
	)

	if err != nil {
		return err
	}

	systemNotifier := system.NewSystemNotifier(m.logger, m.db)

	notifiers := map[string]service.Notifier{
		"smtp":   smtpNotifier,
		"system": systemNotifier,
	}

	m.service = service.NewNotificationService(notifiers, m.logger)

	return nil
}

func (m *NotificationModule) GetNotificationService() *service.NotificationService {
	return m.service
}
