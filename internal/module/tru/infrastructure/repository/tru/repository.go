package tru_repository

import (
	"context"

	tru_entity "github.com/Fi44er/sdmed/internal/module/tru/entity"
	tru_model "github.com/Fi44er/sdmed/internal/module/tru/infrastructure/repository/model"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ITRUCodeRepository interface {
	UpsertMany(ctx context.Context, codes []*tru_entity.TRUCode) error
	GetByID(ctx context.Context, id string) (*tru_entity.TRUCode, error)
	GetByCode(ctx context.Context, code string) (*tru_entity.TRUCode, error)
	GetMany(ctx context.Context, offset, limit int) ([]tru_entity.TRUCode, error)
}

type TRUCodeRepository struct {
	db        *gorm.DB
	logger    *logger.Logger
	converter *Converter
}

func NewTRUCodeRepository(db *gorm.DB, logger *logger.Logger) ITRUCodeRepository {
	return &TRUCodeRepository{
		db:        db,
		logger:    logger,
		converter: &Converter{},
	}
}

func (r *TRUCodeRepository) UpsertMany(ctx context.Context, codes []*tru_entity.TRUCode) error {
	r.logger.Infof("Upserting %d TRU codes", len(codes))
	db := r.db.WithContext(ctx)

	models := make([]*tru_model.TRUCode, 0, len(codes))
	for _, item := range codes {
		models = append(models, r.converter.ToModel(item))
	}

	for _, item := range models {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{"is_custom"}),
		}).Create(item).Error; err != nil {
			return err
		}

		if err := db.Where("tru_code_id = ?", item.ID).Delete(&tru_model.TRUCodePrice{}).Error; err != nil {
			return err
		}

		if len(item.Prices) > 0 {
			for i := range item.Prices {
				item.Prices[i].TRUCodeID = item.ID
			}
			if err := db.Create(&item.Prices).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *TRUCodeRepository) GetByID(ctx context.Context, id string) (*tru_entity.TRUCode, error) {
	var result tru_model.TRUCode
	err := r.db.WithContext(ctx).
		Preload("Prices").
		First(&result, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	entity := r.converter.ToEntity(&result)
	return entity, nil
}

func (r *TRUCodeRepository) GetMany(ctx context.Context, limit, offset int) ([]tru_entity.TRUCode, error) {
	var results []tru_model.TRUCode

	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Preload("Prices").
		Order("code ASC").
		Find(&results).Error

	entities := make([]tru_entity.TRUCode, 0, len(results))
	for _, item := range results {
		entities = append(entities, *r.converter.ToEntity(&item))
	}
	return entities, err
}

func (r *TRUCodeRepository) GetByCode(ctx context.Context, code string) (*tru_entity.TRUCode, error) {
	var result tru_model.TRUCode
	err := r.db.WithContext(ctx).
		Preload("Prices").
		First(&result, "code = ?", code).Error

	if err != nil {
		return nil, err
	}

	entity := r.converter.ToEntity(&result)
	return entity, nil
}
