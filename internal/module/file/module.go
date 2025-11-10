package module

import (
	"time"

	"github.com/Fi44er/sdmed/internal/config"
	file_http "github.com/Fi44er/sdmed/internal/module/file/delivery/http"
	"github.com/Fi44er/sdmed/internal/module/file/infrastucture/filesystem"
	repository "github.com/Fi44er/sdmed/internal/module/file/infrastucture/repository/file"
	file_usecase "github.com/Fi44er/sdmed/internal/module/file/usecase/file"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/postgres/uow"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FileModule struct {
	fileRepository repository.IFileRepository
	fileStorage    *filesystem.LocalFileStorage
	fileUsecase    file_usecase.IFileUsecase
	fileHandler    *file_http.FileHandler
	fileCleaner    *file_usecase.FileCleaner

	logger    *logger.Logger
	validator *validator.Validate
	db        *gorm.DB
	config    *config.Config
	uow       uow.Uow
}

func NewFileModule(
	logger *logger.Logger,
	validator *validator.Validate,
	db *gorm.DB,
	config *config.Config,
	uow uow.Uow,
) *FileModule {
	return &FileModule{
		logger:    logger,
		validator: validator,
		db:        db,
		config:    config,
		uow:       uow,
	}
}

func (m *FileModule) Init() {
	m.uow.RegisterRepository("file", func(tx *gorm.DB) (any, error) {
		return repository.NewFileRepository(m.logger, tx), nil
	})

	m.fileRepository = repository.NewFileRepository(m.logger, m.db)
	m.fileStorage = filesystem.NewLocalFileStorage(m.logger, m.config)
	m.fileUsecase = file_usecase.NewFileUsecase(m.fileRepository, m.uow, m.fileStorage, m.logger, m.config)
	m.fileHandler = file_http.NewFileHandler(m.fileUsecase, m.validator, m.logger)
	m.fileCleaner = file_usecase.NewFileCleaner(m.fileRepository, m.fileStorage, m.logger, 20*time.Minute)
}

func (m *FileModule) InitDelivery(router fiber.Router) {
	m.fileHandler.RegisterRoutes(router)
}

func (m *FileModule) GetFileCleaner() *file_usecase.FileCleaner {
	return m.fileCleaner
}
