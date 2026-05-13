package tru_entity

type TRUCode struct {
	ID       string
	Code     string
	IsCustom bool
	Prices   []TRUCodePrice
}
