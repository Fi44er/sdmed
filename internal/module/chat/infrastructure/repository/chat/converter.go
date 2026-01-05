package chat_repository

import (
	chat_entity "github.com/Fi44er/sdmed/internal/module/chat/entity"
	chat_models "github.com/Fi44er/sdmed/internal/module/chat/infrastructure/repository/model"
)

type Converter struct{}

func (c *Converter) ToModel(entity *chat_entity.Chat) *chat_models.Chat {
	return &chat_models.Chat{
		ID:         entity.ID,
		ClientID:   entity.ClientID,
		OperatorID: entity.OperatorID,
		Subject:    entity.Subject,
		Status:     chat_models.ChatStatus(entity.Status),
		Priority:   entity.Priority,
		Tags:       entity.Tags,

		ClosedAt: entity.ClosedAt,
	}
}

func (c *Converter) ToEntity(model *chat_models.Chat) *chat_entity.Chat {
	return &chat_entity.Chat{
		ID:         model.ID,
		ClientID:   model.ClientID,
		OperatorID: model.OperatorID,
		Subject:    model.Subject,
		Status:     chat_entity.ChatStatus(model.Status),
		Priority:   model.Priority,
		Tags:       model.Tags,

		ClosedAt: model.ClosedAt,
	}
}
