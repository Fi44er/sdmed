package module

import (
	"github.com/Fi44er/sdmed/internal/config"
	"github.com/Fi44er/sdmed/internal/module/file/delivery/http"
	"github.com/Fi44er/sdmed/internal/module/file/infrastucture/filesystem"
	repository "github.com/Fi44er/sdmed/internal/module/file/infrastucture/repository/file"
	usecase "github.com/Fi44er/sdmed/internal/module/file/usecase/file"
	"github.com/Fi44er/sdmed/pkg/logger"
	"github.com/Fi44er/sdmed/pkg/postgres/uow"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type FileModule struct {
	fileRepository *repository.FileRepository
	fileStorage    *filesystem.LocalFileStorage
	fileUsecase    *usecase.FileUsecase
	fileHandler    *http.FileHandler

	logger *logger.Logger
	db     *gorm.DB
	config *config.Config
	uow    uow.Uow
}

func NewFileModule(
	logger *logger.Logger,
	db *gorm.DB,
	config *config.Config,
	uow uow.Uow,
) *FileModule {
	return &FileModule{
		logger: logger,
		db:     db,
		config: config,
		uow:    uow,
	}
}

func (m *FileModule) Init() {
	m.uow.RegisterRepository("file", func(tx *gorm.DB) (interface{}, error) {
		return repository.NewFileRepository(m.logger, tx), nil
	})

	m.fileRepository = repository.NewFileRepository(m.logger, m.db)
	m.fileStorage = filesystem.NewLocalFileStorage(m.logger, m.config)
	m.fileUsecase = usecase.NewFileUsecase(m.fileRepository, m.uow, m.fileStorage, m.logger)
	m.fileHandler = http.NewFileHandler(m.fileUsecase, m.logger, nil)
}

func (m *FileModule) InitDelivery(router fiber.Router) {
	m.fileHandler.RegisterRoutes(router)
}
