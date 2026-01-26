package productapp

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/product"
	"github.com/RubenRodrigo/go-tiny-store/pkg/pagination"
)

// Service handles product-related use cases
type Service struct {
	productRepo product.Repository
}

// NewService creates a new product application service
func NewService(productRepo product.Repository) *Service {
	return &Service{
		productRepo: productRepo,
	}
}

func (s *Service) List(params pagination.Params, filters ProductFilters) (pagination.Result[*product.Product], error) {
	// Convert application filters to domain filters
	domainFilters := product.Filters{
		CategoryID: filters.CategoryID,
		MinPrice:   filters.MinPrice,
		MaxPrice:   filters.MaxPrice,
		Disabled:   filters.Disabled,
	}

	products, count, err := s.productRepo.ListProducts(params, domainFilters)
	if err != nil {
		return pagination.Result[*product.Product]{}, err
	}

	return pagination.BuildResult(params, count, products), nil
}

func (s *Service) Get(id string) (*product.Product, error) {
	return s.productRepo.GetProduct(id)
}

func (s *Service) Create(name string, price float64, stock int, categoryID string) (*product.Product, error) {
	p := &product.Product{
		Name:       name,
		Price:      price,
		Stock:      stock,
		CategoryID: categoryID,
		Disabled:   false,
	}

	err := s.productRepo.CreateProduct(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Service) Update(id, name string, price float64, stock int, disabled bool) (*product.Product, error) {
	p := &product.Product{
		ID:       id,
		Name:     name,
		Price:    price,
		Stock:    stock,
		Disabled: disabled,
	}

	err := s.productRepo.UpdateProduct(p)
	if err != nil {
		return nil, err
	}

	return p, nil
}
