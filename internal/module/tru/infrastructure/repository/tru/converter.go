package tru_repository

import (
	tru_entity "github.com/Fi44er/sdmed/internal/module/tru/entity"
	tru_model "github.com/Fi44er/sdmed/internal/module/tru/infrastructure/repository/model"
)

type Converter struct{}

func (c *Converter) ToModel(e *tru_entity.TRUCode) *tru_model.TRUCode {
	prices := make([]tru_model.TRUCodePrice, 0, len(e.Prices))
	for _, price := range e.Prices {
		prices = append(prices, tru_model.TRUCodePrice{
			ID:        price.ID,
			TRUCodeID: price.TRUCodeID,
			RegionIso: price.RegionIso,
			Price:     price.Price,
		})
	}

	return &tru_model.TRUCode{
		ID:       e.ID,
		Code:     e.Code,
		IsCustom: e.IsCustom,
		Prices:   prices,
	}
}

func (c *Converter) ToEntity(e *tru_model.TRUCode) *tru_entity.TRUCode {
	prices := make([]tru_entity.TRUCodePrice, 0, len(e.Prices))
	for _, price := range e.Prices {
		prices = append(prices, tru_entity.TRUCodePrice{
			ID:        price.ID,
			TRUCodeID: price.TRUCodeID,
			RegionIso: price.RegionIso,
			Price:     price.Price,
		})
	}

	return &tru_entity.TRUCode{
		ID:       e.ID,
		Code:     e.Code,
		IsCustom: e.IsCustom,
		Prices:   prices,
	}
}
