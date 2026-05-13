package tru_usecase

import (
	"context"

	tru_entity "github.com/Fi44er/sdmed/internal/module/tru/entity"
	"github.com/Fi44er/sdmed/pkg/logger"
)

type ITRUCodeRepository interface {
	UpsertMany(ctx context.Context, codes []*tru_entity.TRUCode) error
	GetByID(ctx context.Context, id string) (*tru_entity.TRUCode, error)
	GetByCode(ctx context.Context, code string) (*tru_entity.TRUCode, error)
	GetMany(ctx context.Context, offset, limit int) ([]tru_entity.TRUCode, error)
}

type ITRUCodeUseCase interface {
	UpsertMany(ctx context.Context, codes []*tru_entity.TRUCode) error
	GetByID(ctx context.Context, id string) (*tru_entity.TRUCode, error)
	GetByCode(ctx context.Context, code string) (*tru_entity.TRUCode, error)
	GetMany(ctx context.Context, offset, limit int) ([]tru_entity.TRUCode, error)
}

type TRUCodeUseCase struct {
	logger *logger.Logger
	repo   ITRUCodeRepository
}

func NewTRUCodeUseCase(logger *logger.Logger, repo ITRUCodeRepository) ITRUCodeUseCase {
	return &TRUCodeUseCase{
		logger: logger,
		repo:   repo,
	}
}

func (u *TRUCodeUseCase) UpsertMany(ctx context.Context, codes []*tru_entity.TRUCode) error {
	if len(codes) == 0 {
		return nil
	}
	return u.repo.UpsertMany(ctx, codes)
}

func (u *TRUCodeUseCase) GetByID(ctx context.Context, id string) (*tru_entity.TRUCode, error) {
	return u.repo.GetByID(ctx, id)
}

func (u *TRUCodeUseCase) GetByCode(ctx context.Context, code string) (*tru_entity.TRUCode, error) {
	return u.repo.GetByCode(ctx, code)
}

func (u *TRUCodeUseCase) GetMany(ctx context.Context, offset, limit int) ([]tru_entity.TRUCode, error) {
	return u.repo.GetMany(ctx, offset, limit)
}
