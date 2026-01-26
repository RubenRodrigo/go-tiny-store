package product

import "github.com/RubenRodrigo/go-tiny-store/pkg/pagination"

// Repository defines the interface for product persistence operations
type Repository interface {
	ListProducts(params pagination.Params, filters Filters) ([]*Product, int64, error)
	GetProduct(id string) (*Product, error)
	CreateProduct(product *Product) error
	UpdateProduct(product *Product) error
}
