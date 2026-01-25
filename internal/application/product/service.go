package product

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/models"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/repository"
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
func (p *Service) List() ([]*models.Product, error) {
	panic("unimplemented")
}

// Update implements ProductService.
func (p *Service) Update(name string, id string) (*models.Product, error) {
	panic("unimplemented")
}
