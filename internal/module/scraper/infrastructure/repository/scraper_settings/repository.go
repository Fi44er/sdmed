package scraper_settings_repository

import (
	"context"

	scraper_entity "github.com/Fi44er/sdmed/internal/module/scraper/entity"
	scraper_models "github.com/Fi44er/sdmed/internal/module/scraper/infrastructure/repository/models"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type IScraperSettingsRepository interface {
	Update(ctx context.Context, settings *scraper_entity.ScraperSettings) error
	Get(ctx context.Context) (*scraper_entity.ScraperSettings, error)
}

type ScraperSettingsRepository struct {
	logger    *logger.Logger
	db        *gorm.DB
	converter *Converter
}

func NewScraperSettingsRepository(db *gorm.DB, logger *logger.Logger) IScraperSettingsRepository {
	return &ScraperSettingsRepository{
		logger:    logger,
		db:        db,
		converter: &Converter{},
	}
}

func (r *ScraperSettingsRepository) Update(ctx context.Context, settings *scraper_entity.ScraperSettings) error {
	r.logger.Infof("updating scraper settings: %v", settings)
	model := r.converter.ToModel(settings)
	r.logger.Warnf("Model: %+v", model)
	model.ID = 1
	if err := r.db.WithContext(ctx).Save(model).Error; err != nil {
		r.logger.Errorf("failed to update scraper settings: %v", err)
		return err
	}
	r.logger.Info("scraper settings updated successfully")
	return nil
}

func (r *ScraperSettingsRepository) Get(ctx context.Context) (*scraper_entity.ScraperSettings, error) {
	r.logger.Info("getting scraper settings")
	var model scraper_models.ScraperSettings
	if err := r.db.WithContext(ctx).First(&model).Error; err != nil {
		r.logger.Errorf("failed to get scraper settings: %v", err)
		return nil, err
	}
	settings := r.converter.ToEntity(&model)
	r.logger.Info("scraper settings retrieved successfully")
	return settings, nil
}
