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

// Create implements CategoryService.
func (s *categoryService) Save(name string, id uint) (*models.Category, error) {
	category := &models.Category{Name: name, ID: id}

	if err := s.categoryRepo.Save(category); err != nil {
		return nil, err
	}

	return category, nil
}

// Delete implements CategoryService.
func (s *categoryService) Delete(id uint) (*models.Category, error) {
	category := &models.Category{ID: id}

	if err := s.categoryRepo.Delete(category); err != nil {
		return nil, err
	}

	return category, nil
}

// List implements CategoryService.
func (s *categoryService) List() ([]*models.Category, error) {
	categories, err := s.categoryRepo.List()
	if err != nil {
		return nil, err
	}

	return categories, nil
}
