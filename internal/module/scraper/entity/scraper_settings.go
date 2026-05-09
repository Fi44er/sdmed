package scraper_entity

import (
	"time"
)

type ScraperSettings struct {
	ID        uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	IntervalMs    time.Duration
	MaxRetries    int
	RetryDelayS   time.Duration
	PauseS        time.Duration
	TimeoutS      time.Duration
	MaxGoroutines int

	ConfigName string
}
