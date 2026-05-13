package tru_module

import (
	tru_repository "github.com/Fi44er/sdmed/internal/module/tru/infrastructure/repository/tru"
	tru_usecase "github.com/Fi44er/sdmed/internal/module/tru/usecase"
	"github.com/Fi44er/sdmed/pkg/logger"
	"gorm.io/gorm"
)

type TruModule struct {
	truRepository tru_repository.ITRUCodeRepository
	truUseCase    tru_usecase.ITRUCodeUseCase

	logger *logger.Logger
	db     *gorm.DB
}

func NewTruModule(logger *logger.Logger, db *gorm.DB) *TruModule {
	return &TruModule{
		logger: logger,
		db:     db,
	}
}

func (m *TruModule) Init() {
	m.truRepository = tru_repository.NewTRUCodeRepository(m.db, m.logger)
	m.truUseCase = tru_usecase.NewTRUCodeUseCase(m.logger, m.truRepository)
}

func (m *TruModule) GetTRUCodeUsecase() tru_usecase.ITRUCodeUseCase {
	return m.truUseCase
}
