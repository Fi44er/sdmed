package module

import (
	"github.com/Fi44er/sdmed/internal/config"
	"github.com/Fi44er/sdmed/internal/module/file/infrastucture/filesystem"
	repository "github.com/Fi44er/sdmed/internal/module/file/infrastucture/repository/file"
	usecase "github.com/Fi44er/sdmed/internal/module/file/usecase/file"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type FileModule struct {
	fileRepository *repository.FileRepository
	fileStorage    *filesystem.LocalFileStorage
	fileUsecase    *usecase.FileUsecase

	logger *logger.Logger
	db     *gorm.DB
	config *config.Config
}
