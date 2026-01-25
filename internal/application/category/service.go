package category

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/models"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/repository"
)

type Service struct {
	repository repository.Category
}

func NewService(repository repository.Category) *Service {
	return &Service{
		repository: repository,
	}
}

// List implements CategoryService.
func (s *Service) List() ([]*models.Category, error) {
	categories, err := s.repository.ListCategory()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// Create implements CategoryService.
func (s *Service) Create(name string) (*models.Category, error) {
	category := &models.Category{Name: name}

	if err := s.repository.CreateCategory(category); err != nil {
		return nil, err
	}

	return category, nil
}

// Create implements CategoryService.
func (s *Service) Update(name, id string) (*models.Category, error) {
	category := &models.Category{Name: name}

	if err := s.repository.UpdateCategory(id, category); err != nil {
		return nil, err
	}

	return category, nil
}

// Delete implements CategoryService.
func (s *Service) Delete(id string) error {
	if err := s.repository.DeleteCategory(id); err != nil {
		return err
	}

	return nil
}
