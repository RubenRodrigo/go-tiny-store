package repository

import (
	"errors"
	"log"

	"github.com/RubenRodrigo/go-tiny-store/internal/apperrors"
	"github.com/RubenRodrigo/go-tiny-store/internal/models"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

// Delete implements CategoryRepository.
func (r *categoryRepository) Delete(category *models.Category) error {
	err := r.db.Delete(category).Error

	if err == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound
		}

		log.Printf("ERROR: Failed to delete category in database. Name: %s, Error: %v",
			category.Name, err)

		return apperrors.ErrDatabaseError
	}

	return nil
}

// GetById implements CategoryRepository.
func (r *categoryRepository) GetById(id string) (*models.Category, error) {
	var category models.Category
	if err := r.db.Where("id = ?", id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}

		log.Printf("ERROR: Failed to read user in database. Id: %s, Error: %v",
			id, err)

		return nil, apperrors.ErrDatabaseError
	}

	return &category, nil
}

// List implements CategoryRepository.
func (r *categoryRepository) List() ([]*models.Category, error) {
	var categories []*models.Category
	if err := r.db.Find(&categories).Error; err != nil {
		log.Printf("ERROR: Failed to read categories in database. Error: %s", err)
		return nil, apperrors.ErrDatabaseError
	}

	return categories, nil
}

func (r *categoryRepository) Save(category *models.Category) error {
	err := r.db.Save(category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperrors.ErrDuplicateEntry
		}

		log.Printf("ERROR: Failed to update category in database. ID: %d, Name: %s, Error: %v",
			category.ID, category.Name, err)

		return apperrors.ErrDatabaseError
	}
	return nil
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}
