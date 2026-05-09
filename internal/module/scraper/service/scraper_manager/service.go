package scrapermanager_service

import "github.com/Fi44er/sdmed/pkg/logger"

type IScraperManagerService interface {
}

type ScraperManagerService struct {
	logger *logger.Logger
}
