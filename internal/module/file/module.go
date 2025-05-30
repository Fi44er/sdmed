package module

import (
	"github.com/Fi44er/sdmedik/backend/internal/config"
	"github.com/Fi44er/sdmedik/backend/internal/module/file/infrastucture/filesystem"
	repository "github.com/Fi44er/sdmedik/backend/internal/module/file/infrastucture/repository/file"
	usecase "github.com/Fi44er/sdmedik/backend/internal/module/file/usecase/file"
	"github.com/Fi44er/sdmedik/backend/pkg/logger"
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
