package product_entity

import (
	"time"

	"github.com/Fi44er/sdmed/pkg/utils"
)

type Category struct {
	ID        string
	Name      string
	Slug      string
	Images    []File
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time

	Characteristics []Characteristic
}

func (c *Category) Slugify() {
	c.Slug = utils.CreateSlugRU(c.Name)
}
