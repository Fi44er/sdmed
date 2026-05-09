package scraper_dto

type UpdateSettingsDTO struct {
	IntervalMs    int `json:"interval_ms" validate:"required,min=500"`
	MaxRetries    int `json:"max_retries" validate:"required,min=1,max=10"`
	RetryDelayS   int `json:"retry_delay_s" validate:"required,min=1"`
	PauseS        int `json:"pause_s" validate:"required,min=1"`
	TimeoutS      int `json:"timeout_s" validate:"required,min=1"`
	MaxGoroutines int `json:"max_goroutines" validate:"required,min=1,max=10"`
}
