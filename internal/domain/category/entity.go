package category

import "time"

// Category represents a product category (pure domain entity)
type Category struct {
	ID        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
