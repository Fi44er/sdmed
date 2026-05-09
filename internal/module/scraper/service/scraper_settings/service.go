package scrapersettings_service

import (
	"context"

	scraper_entity "github.com/Fi44er/sdmed/internal/module/scraper/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
)

type IScraperSettingsRepository interface {
	Update(ctx context.Context, settings *scraper_entity.ScraperSettings) error
	Get(ctx context.Context) (*scraper_entity.ScraperSettings, error)
}

type IScraperSettingsService interface {
	Update(ctx context.Context, settings *scraper_entity.ScraperSettings) error
	Get(ctx context.Context) (*scraper_entity.ScraperSettings, error)
}

type ScraperSettingsService struct {
	logger     *logger.Logger
	repository IScraperSettingsRepository
}

func NewScraperSettingsService(logger *logger.Logger, repository IScraperSettingsRepository) IScraperSettingsService {
	return &ScraperSettingsService{
		logger:     logger,
		repository: repository,
	}
}

func (s *ScraperSettingsService) Update(ctx context.Context, settings *scraper_entity.ScraperSettings) error {
	s.logger.Warnf("updating scraper settings: %+v", settings)
	if err := s.repository.Update(ctx, settings); err != nil {
		return err
	}
	return nil
}

func (s *ScraperSettingsService) Get(ctx context.Context) (*scraper_entity.ScraperSettings, error) {
	settings, err := s.repository.Get(ctx)
	if err != nil {
		return nil, err
	}
	return settings, nil
}
