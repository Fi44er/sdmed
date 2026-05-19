package cart_adapters

import (
	"context"

	tru_entity "github.com/Fi44er/sdmed/internal/module/tru/entity"
	tru_usecase "github.com/Fi44er/sdmed/internal/module/tru/usecase"
)

type TRUCodeUsecaseAdapter struct {
	truUseCase tru_usecase.ITRUCodeUseCase
}

func NewTRUCodeUsecaseAdapter(truUseCase tru_usecase.ITRUCodeUseCase) *TRUCodeUsecaseAdapter {
	return &TRUCodeUsecaseAdapter{
		truUseCase: truUseCase,
	}
}

func (a *TRUCodeUsecaseAdapter) GetByID(ctx context.Context, id string) (*tru_entity.TRUCode, error) {
	return a.truUseCase.GetByID(ctx, id)
}
