package product

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/models"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/repository"
	"github.com/RubenRodrigo/go-tiny-store/pkg/pagination"
)

type Service struct {
	productRepo repository.Product
}

func NewService(productRepo repository.Product) *Service {
	return &Service{
		productRepo: productRepo,
	}
}

// Create implements ProductService.
func (p *Service) Create(name string) (*models.Product, error) {
	panic("unimplemented")
}

// Delete implements ProductService.
func (p *Service) Delete(id string) error {
	panic("unimplemented")
}

// List implements ProductService.
func (p *Service) List(params pagination.Params) (pagination.Result[*models.Product], error) {
	products, count, err := p.productRepo.ListProducts(params)
	if err != nil {
		return pagination.Result[*models.Product]{}, err
	}

	return pagination.BuildResult(params, count, products), nil
}

func (p *Service) Get(id string) (*models.Product, error) {
	product, err := p.productRepo.GetProduct(id)
	if err != nil {
		return nil, err
	}

	return product, nil
}

// Update implements ProductService.
func (p *Service) Update(name string, id string) (*models.Product, error) {
	panic("unimplemented")
}
