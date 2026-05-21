package system

import (
	notification_models "github.com/Fi44er/sdmed/internal/module/notification/infrastructure/repository/model"
	"github.com/Fi44er/sdmed/internal/module/notification/service"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type SystemNotifier struct {
	logger *logger.Logger
	db     *gorm.DB
}

func NewSystemNotifier(logger *logger.Logger, db *gorm.DB) *SystemNotifier {
	return &SystemNotifier{
		logger: logger,
		db:     db,
	}
}

func (n *SystemNotifier) Send(msg *service.Message) {
	n.logger.Infof("[SYSTEM NOTIFICATION] Saving to DB - Recipient: %s | Subject: %s", msg.Recipient, msg.Subject)

	notif := &notification_models.Notification{
		UserID:  msg.Recipient,
		Subject: msg.Subject,
		Content: msg.Content,
		IsRead:  false,
	}

	if err := n.db.Create(notif).Error; err != nil {
		n.logger.Errorf("Failed to save system notification to DB: %v", err)
	}
}
