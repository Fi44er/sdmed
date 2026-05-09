package scraper_settings_repository

import (
	"time"

	scraper_entity "github.com/Fi44er/sdmed/internal/module/scraper/entity"
	scraper_models "github.com/Fi44er/sdmed/internal/module/scraper/infrastructure/repository/models"
	"gorm.io/gorm"
)

type Converter struct {
}

func (c *Converter) ToEntity(m *scraper_models.ScraperSettings) *scraper_entity.ScraperSettings {
	if m == nil {
		return nil
	}

	return &scraper_entity.ScraperSettings{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		DeletedAt: m.DeletedAt.Time,

		IntervalMs:    time.Duration(m.IntervalMs) * time.Millisecond,
		MaxRetries:    m.MaxRetries,
		RetryDelayS:   time.Duration(m.RetryDelayS) * time.Second,
		PauseS:        time.Duration(m.PauseS) * time.Second,
		TimeoutS:      time.Duration(m.TimeoutS) * time.Second,
		MaxGoroutines: m.MaxGoroutines,

		ConfigName: m.ConfigName,
	}
}

func (c *Converter) ToModel(e *scraper_entity.ScraperSettings) *scraper_models.ScraperSettings {
	if e == nil {
		return nil
	}

	return &scraper_models.ScraperSettings{
		ID:        e.ID,
		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
		DeletedAt: gorm.DeletedAt{Time: e.DeletedAt, Valid: !e.DeletedAt.IsZero()},

		IntervalMs:    int(e.IntervalMs.Milliseconds()),
		MaxRetries:    e.MaxRetries,
		RetryDelayS:   int(e.RetryDelayS.Seconds()),
		PauseS:        int(e.PauseS.Seconds()),
		TimeoutS:      int(e.TimeoutS.Seconds()),
		MaxGoroutines: e.MaxGoroutines,

		ConfigName: e.ConfigName,
	}
}
