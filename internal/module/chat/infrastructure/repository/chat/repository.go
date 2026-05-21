package chat_repository

import (
	"context"

	chat_entity "github.com/Fi44er/sdmed/internal/module/chat/entity"
	chat_models "github.com/Fi44er/sdmed/internal/module/chat/infrastructure/repository/model"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/utils"
	"gorm.io/gorm"
)

type IChatRepository interface {
	Create(ctx context.Context, chat *chat_entity.Chat) error
	GetAll(ctx context.Context, page, pageSize int) ([]chat_entity.Chat, int64, error)
	GetByID(ctx context.Context, id string) (*chat_entity.Chat, error)
	GetByClientID(ctx context.Context, clientID string, page, pageSize int) ([]chat_entity.Chat, int64, error)
	GetByOperatorID(ctx context.Context, operatorID string, page, pageSize int) ([]chat_entity.Chat, int64, error)
	Update(ctx context.Context, chat *chat_entity.Chat) error
}

type ChatRepository struct {
	logger    *logger.Logger
	db        *gorm.DB
	converter *Converter
}

func NewChatRepository(logger *logger.Logger, db *gorm.DB) IChatRepository {
	return &ChatRepository{
		logger:    logger,
		db:        db,
		converter: &Converter{},
	}
}

func (r *ChatRepository) Create(ctx context.Context, chat *chat_entity.Chat) error {
	r.logger.Infof("Creating chat: %s", chat.Subject)
	chatModel := r.converter.ToModel(chat)
	err := r.db.WithContext(ctx).Create(&chatModel).Error
	if err != nil {
		r.logger.Errorf("Failed to create chat: %v", err)
		return err
	}
	chat.ID = chatModel.ID
	chat.CreatedAt = chatModel.CreatedAt
	chat.UpdatedAt = chatModel.UpdatedAt
	r.logger.Infof("Created chat successfully: %s", chat.ID)
	return nil
}

func (r *ChatRepository) GetAll(ctx context.Context, page, pageSize int) ([]chat_entity.Chat, int64, error) {
	r.logger.Debug("Getting chats")
	var (
		chatsModels []chat_models.Chat
		total       int64
	)

	query := r.db.WithContext(ctx).Model(&chat_models.Chat{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := utils.SafeCalculateForPostgres(page, pageSize)

	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}

	err := query.
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&chatsModels).Error

	if err != nil {
		r.logger.Errorf("Failed to get chats: %v", err)
		return nil, 0, err
	}

	chats := make([]chat_entity.Chat, len(chatsModels))
	for i, chatModel := range chatsModels {
		chats[i] = *r.converter.ToEntity(&chatModel)
	}

	return chats, total, nil
}

func (r *ChatRepository) GetByID(ctx context.Context, id string) (*chat_entity.Chat, error) {
	r.logger.Infof("Getting chat by ID: %s", id)
	var chatModel chat_models.Chat
	err := r.db.WithContext(ctx).First(&chatModel, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.logger.Errorf("Failed to get chat by ID %s: %v", id, err)
		return nil, err
	}
	return r.converter.ToEntity(&chatModel), nil
}

func (r *ChatRepository) GetByClientID(ctx context.Context, clientID string, page, pageSize int) ([]chat_entity.Chat, int64, error) {
	r.logger.Infof("Getting chats for client: %s", clientID)
	var (
		chatsModels []chat_models.Chat
		total       int64
	)

	query := r.db.WithContext(ctx).Model(&chat_models.Chat{}).Where("client_id = ?", clientID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := utils.SafeCalculateForPostgres(page, pageSize)

	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}

	err := query.
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&chatsModels).Error

	if err != nil {
		r.logger.Errorf("Failed to get chats for client: %v", err)
		return nil, 0, err
	}

	chats := make([]chat_entity.Chat, len(chatsModels))
	for i, chatModel := range chatsModels {
		chats[i] = *r.converter.ToEntity(&chatModel)
	}

	return chats, total, nil
}

func (r *ChatRepository) GetByOperatorID(ctx context.Context, operatorID string, page, pageSize int) ([]chat_entity.Chat, int64, error) {
	r.logger.Infof("Getting chats for operator: %s", operatorID)
	var (
		chatsModels []chat_models.Chat
		total       int64
	)

	query := r.db.WithContext(ctx).Model(&chat_models.Chat{}).Where("operator_id = ? OR operator_id IS NULL", operatorID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := utils.SafeCalculateForPostgres(page, pageSize)

	if limit <= 0 {
		limit = -1
	}
	if offset <= 0 {
		offset = -1
	}

	err := query.
		Order("updated_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&chatsModels).Error

	if err != nil {
		r.logger.Errorf("Failed to get chats for operator: %v", err)
		return nil, 0, err
	}

	chats := make([]chat_entity.Chat, len(chatsModels))
	for i, chatModel := range chatsModels {
		chats[i] = *r.converter.ToEntity(&chatModel)
	}

	return chats, total, nil
}

func (r *ChatRepository) Update(ctx context.Context, chat *chat_entity.Chat) error {
	r.logger.Infof("Updating chat: %s", chat.ID)
	chatModel := r.converter.ToModel(chat)
	err := r.db.WithContext(ctx).Model(&chat_models.Chat{}).Where("id = ?", chat.ID).Updates(chatModel).Error
	if err != nil {
		r.logger.Errorf("Failed to update chat %s: %v", chat.ID, err)
		return err
	}
	return nil
}
