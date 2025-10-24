package product_entity

import (
	"fmt"
	"regexp"
	"strings"
)

type Price struct {
	Price     float64
	TRU       string
	RegionISO string
}

type Product struct {
	ID          string
	Article     string
	Name        string
	Description string
	Price       Price
	Images      []string

	CategoryID      string
	Characteristics []CharacteristicValue
}

func (p *Product) Validate() error {
	switch {
	case p.ValidateArticle() != nil:
		return fmt.Errorf("invalid article")
	case p.ValidateName() != nil:
		return fmt.Errorf("invalid name")
	case p.ValidatePrice() != nil:
		return fmt.Errorf("invalid price")
	}
	return nil
}

func (p *Product) ValidateArticle() error {
	if strings.TrimSpace(p.Article) == "" {
		return fmt.Errorf("article cannot be empty")
	}
	articleRegex := regexp.MustCompile(`^[A-Z]{3}-\d{3}$`)
	if !articleRegex.MatchString(p.Article) {
		return fmt.Errorf("article must match pattern 'ABC-123'")
	}
	return nil
}

func (p *Product) ValidateName() error {
	name := strings.TrimSpace(p.Name)
	if name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("name is too long")
	}
	return nil
}

func (p *Product) ValidatePrice() error {
	if p.Price.Price <= 0 {
		return fmt.Errorf("price must be positive")
	}
	return nil
}

func (p *Product) FormatDescription() {
	p.Description = strings.TrimSpace(p.Description)
}

func (p *Product) IsInCategory(categoryID string) bool {
	return p.CategoryID == categoryID
}

func (p *Product) ApplyDiscount(discountPercent float64) error {
	if discountPercent < 0 || discountPercent > 100 {
		return fmt.Errorf("invalid discount percentage")
	}
	p.Price.Price = p.Price.Price * (1 - discountPercent/100)
	return nil
}

func (p *Product) IncreasePrice(percent float64) error {
	if percent < 0 {
		return fmt.Errorf("invalid percentage")
	}
	p.Price.Price = p.Price.Price * (1 + percent/100)
	return nil
}
