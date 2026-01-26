package categoryapp

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/category"
)

// Service handles category-related use cases
type Service struct {
	categoryRepo category.Repository
}

// NewService creates a new category application service
func NewService(categoryRepo category.Repository) *Service {
	return &Service{
		categoryRepo: categoryRepo,
	}
}

func (s *Service) List() ([]*category.Category, error) {
	return s.categoryRepo.ListCategories()
}

func (s *Service) GetByID(id string) (*category.Category, error) {
	return s.categoryRepo.GetCategoryByID(id)
}

func (s *Service) Create(name string) (*category.Category, error) {
	c := &category.Category{Name: name}

	if err := s.categoryRepo.CreateCategory(c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *Service) Update(id, name string) (*category.Category, error) {
	c := &category.Category{Name: name}

	if err := s.categoryRepo.UpdateCategory(id, c); err != nil {
		return nil, err
	}

	return c, nil
}

func (s *Service) Delete(id string) error {
	return s.categoryRepo.DeleteCategory(id)
}
