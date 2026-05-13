package scraper_adapters

import (
	"context"

	scraper_entity "github.com/Fi44er/sdmed/internal/module/scraper/entity"
	tru_entity "github.com/Fi44er/sdmed/internal/module/tru/entity"
)

type ITRUUseCase interface {
	UpsertMany(ctx context.Context, codes []*tru_entity.TRUCode) error
	GetByCode(ctx context.Context, code string) (*tru_entity.TRUCode, error)
}

type TRUUseCaseAdapter struct {
	truUseCase ITRUUseCase
}

func NewTRUUseCaseAdapter(truUseCase ITRUUseCase) *TRUUseCaseAdapter {
	return &TRUUseCaseAdapter{
		truUseCase: truUseCase,
	}
}

func (a *TRUUseCaseAdapter) UpsertMany(ctx context.Context, codes []*scraper_entity.TRUCode) error {
	truCodes := make([]*tru_entity.TRUCode, 0, len(codes))
	for _, code := range codes {
		truCodes = append(truCodes, a.toTRUEntity(code))
	}
	return a.truUseCase.UpsertMany(ctx, truCodes)
}

func (a *TRUUseCaseAdapter) GetByCode(ctx context.Context, code string) (*scraper_entity.TRUCode, error) {
	truCode, err := a.truUseCase.GetByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return a.toScraperTRU(truCode), nil
}

func (a *TRUUseCaseAdapter) toTRUEntity(code *scraper_entity.TRUCode) *tru_entity.TRUCode {
	prices := make([]tru_entity.TRUCodePrice, 0, len(code.Prices))
	for _, price := range code.Prices {
		prices = append(prices, tru_entity.TRUCodePrice{
			ID:        price.ID,
			TRUCodeID: price.TRUCodeID,
			RegionIso: price.RegionIso,
			Price:     price.Price,
		})
	}
	return &tru_entity.TRUCode{
		ID:       code.ID,
		Code:     code.Code,
		IsCustom: code.IsCustom,
		Prices:   prices,
	}
}

func (a *TRUUseCaseAdapter) toScraperTRU(code *tru_entity.TRUCode) *scraper_entity.TRUCode {
	prices := make([]scraper_entity.TRUCodePrice, 0, len(code.Prices))
	for _, price := range code.Prices {
		prices = append(prices, scraper_entity.TRUCodePrice{
			ID:        price.ID,
			TRUCodeID: price.TRUCodeID,
			RegionIso: price.RegionIso,
			Price:     price.Price,
		})
	}
	return &scraper_entity.TRUCode{
		ID:       code.ID,
		Code:     code.Code,
		IsCustom: code.IsCustom,
		Prices:   prices,
	}
}
