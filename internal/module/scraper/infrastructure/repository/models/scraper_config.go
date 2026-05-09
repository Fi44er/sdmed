package scraper_models

import (
	"time"

	"gorm.io/gorm"
)

type ScraperSettings struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	IntervalMs    int `gorm:"column:interval_ms;not null;default:1500"`
	MaxRetries    int `gorm:"column:max_retries;not null;default:3"`
	RetryDelayS   int `gorm:"column:retry_delay_s;not null;default:5"`
	PauseS        int `gorm:"column:pause_s;not null;default:60"`
	TimeoutS      int `gorm:"column:timeout_s;not null;default:30"`
	MaxGoroutines int `gorm:"column:max_goroutines;not null;default:2"`

	ConfigName string `gorm:"column:config_name;default:'Default SFR Parser'"`
}
