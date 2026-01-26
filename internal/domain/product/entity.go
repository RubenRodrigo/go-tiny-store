package product

import "time"

// Product represents a product in the system (pure domain entity)
type Product struct {
	ID         string
	Name       string
	Price      float64
	Disabled   bool
	Stock      int
	CategoryID string
	Images     []ProductImage
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

// ProductImage represents an image associated with a product
type ProductImage struct {
	ID        string
	ProductID string
	URL       string
	AltText   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
