package product_entity

import "time"

type Category struct {
	ID        string
	Name      string
	Images    []File
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Characteristic []Characteristic
}
