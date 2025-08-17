package service

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/models"
	"github.com/RubenRodrigo/go-tiny-store/internal/repository"
)

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{
		categoryRepo: categoryRepo,
	}
}

// List implements CategoryService.
func (s *categoryService) List() ([]*models.Category, error) {
	categories, err := s.categoryRepo.List()
	if err != nil {
		return nil, err
	}

	return categories, nil
}

// Create implements CategoryService.
func (s *categoryService) Create(name string) (*models.Category, error) {
	category := &models.Category{Name: name}

	if err := s.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

// Create implements CategoryService.
func (s *categoryService) Update(name, id string) (*models.Category, error) {
	category := &models.Category{Name: name}

	if err := s.categoryRepo.Update(id, category); err != nil {
		return nil, err
	}

	return category, nil
}

// Delete implements CategoryService.
func (s *categoryService) Delete(id string) error {

	if err := s.categoryRepo.Delete(id); err != nil {
		return err
	}

	return nil
}
