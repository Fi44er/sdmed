package scraper_entity

type TRUCode struct {
	ID       string
	Code     string
	IsCustom bool
	Prices   []TRUCodePrice
}

type TRUCodePrice struct {
	ID        string
	TRUCodeID string
	RegionIso string
	Price     float64
}
