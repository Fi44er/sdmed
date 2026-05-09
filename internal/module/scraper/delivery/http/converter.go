package scraper_http

import (
	"time"

	scraper_dto "github.com/Fi44er/sdmed/internal/module/scraper/dto"
	scraper_entity "github.com/Fi44er/sdmed/internal/module/scraper/entity"
)

type Converter struct {
}

func (c *Converter) ToEntity(dto *scraper_dto.UpdateSettingsDTO) *scraper_entity.ScraperSettings {
	return &scraper_entity.ScraperSettings{
		IntervalMs:    time.Duration(dto.IntervalMs) * time.Millisecond,
		MaxRetries:    dto.MaxRetries,
		RetryDelayS:   time.Duration(dto.RetryDelayS) * time.Second,
		PauseS:        time.Duration(dto.PauseS) * time.Second,
		TimeoutS:      time.Duration(dto.TimeoutS) * time.Second,
		MaxGoroutines: dto.MaxGoroutines,
	}
}
