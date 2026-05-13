package scraper_entity

type Category struct {
	URL  string
	Name string
}

type Item struct {
	Price  float64
	Region string
}

type Items struct {
	CategoryArticle string
	CategoryName    string
	Items           []Item
	Product         []ParseProductsArticlesType
}

type ParseProductsArticlesType struct {
	Article string
	Name    string
}

type ScrapeParams struct {
	Categories []string // Фильтр по разделам (06, 07...)
	Articles   []string // Фильтр по конкретным подкатегориям (06-01-01...)
	Regions    []string // Опционально: фильтр по регионам (ISO коды)
}
