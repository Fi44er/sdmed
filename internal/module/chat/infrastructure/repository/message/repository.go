package message_repository

import (
	"context"
	"time"

	chat_entity "github.com/Fi44er/sdmed/internal/module/chat/entity"
	chat_models "github.com/Fi44er/sdmed/internal/module/chat/infrastructure/repository/model"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/utils"
	"gorm.io/gorm"
)

type IMessageRepository interface {
	Create(ctx context.Context, msg *chat_entity.Message) error
	GetByChatID(ctx context.Context, chatID string, page, pageSize int) ([]chat_entity.Message, error)
	MarkAsRead(ctx context.Context, msgID string) error
	GetByID(ctx context.Context, id string) (*chat_entity.Message, error)
}

type MessageRepository struct {
	logger *logger.Logger
	db     *gorm.DB
}

func NewMessageRepository(logger *logger.Logger, db *gorm.DB) IMessageRepository {
	return &MessageRepository{
		logger: logger,
		db:     db,
	}
}

func (r *MessageRepository) ToModel(entity *chat_entity.Message) *chat_models.Message {
	var replyToID *string
	if entity.ReplyToID != "" {
		replyToID = &entity.ReplyToID
	}
	return &chat_models.Message{
		ID:        entity.ID,
		ChatID:    entity.ChatID,
		SenderID:  entity.SenderID,
		Type:      chat_models.MessageType(entity.Type),
		Payload:   entity.Payload,
		ReplyToID: replyToID,
		IsEdited:  entity.IsEdited,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		ReadAt:    entity.ReadAt,
	}
}

func (r *MessageRepository) ToEntity(model *chat_models.Message) *chat_entity.Message {
	var replyToID string
	if model.ReplyToID != nil {
		replyToID = *model.ReplyToID
	}
	return &chat_entity.Message{
		ID:        model.ID,
		ChatID:    model.ChatID,
		SenderID:  model.SenderID,
		Type:      chat_entity.MessageType(model.Type),
		Payload:   model.Payload,
		ReplyToID: replyToID,
		IsEdited:  model.IsEdited,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
		ReadAt:    model.ReadAt,
	}
}

func (r *MessageRepository) Create(ctx context.Context, msg *chat_entity.Message) error {
	r.logger.Infof("Creating message for chat: %s", msg.ChatID)
	model := r.ToModel(msg)
	if model.CreatedAt.IsZero() {
		model.CreatedAt = time.Now()
	}
	if model.UpdatedAt.IsZero() {
		model.UpdatedAt = time.Now()
	}
	err := r.db.WithContext(ctx).Create(model).Error
	if err != nil {
		r.logger.Errorf("Failed to create message: %v", err)
		return err
	}
	msg.ID = model.ID
	msg.CreatedAt = model.CreatedAt
	msg.UpdatedAt = model.UpdatedAt
	r.logger.Infof("Message created successfully: %s", msg.ID)
	return nil
}

func (r *MessageRepository) GetByChatID(ctx context.Context, chatID string, page, pageSize int) ([]chat_entity.Message, error) {
	r.logger.Infof("Getting messages for chat: %s", chatID)
	var models []chat_models.Message
	offset, limit := utils.SafeCalculateForPostgres(page, pageSize)
	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}

	err := r.db.WithContext(ctx).
		Where("chat_id = ?", chatID).
		Order("created_at ASC").
		Offset(offset).
		Limit(limit).
		Find(&models).Error
	if err != nil {
		r.logger.Errorf("Failed to get messages: %v", err)
		return nil, err
	}

	entities := make([]chat_entity.Message, len(models))
	for i, m := range models {
		entities[i] = *r.ToEntity(&m)
	}
	return entities, nil
}

func (r *MessageRepository) MarkAsRead(ctx context.Context, msgID string) error {
	r.logger.Infof("Marking message as read: %s", msgID)
	now := time.Now()
	err := r.db.WithContext(ctx).
		Model(&chat_models.Message{}).
		Where("id = ? AND read_at IS NULL", msgID).
		Update("read_at", now).Error
	if err != nil {
		r.logger.Errorf("Failed to mark message as read: %v", err)
		return err
	}
	return nil
}

func (r *MessageRepository) GetByID(ctx context.Context, id string) (*chat_entity.Message, error) {
	r.logger.Infof("Getting message by ID: %s", id)
	var model chat_models.Message
	err := r.db.WithContext(ctx).First(&model, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Errorf("Failed to get message: %v", err)
		return nil, err
	}
	return r.ToEntity(&model), nil
}
