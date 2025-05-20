package entity

type Product struct {
	ID          string
	Article     string
	Name        string
	Description string
	Price       float64

	Category []Category
}
