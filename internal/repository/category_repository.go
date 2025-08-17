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
func (r *categoryRepository) Delete(id string) error {
	err := r.db.Delete(&models.Category{}, id).Error

	if err == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound
		}

		log.Printf("ERROR: Failed to delete category in database. ID: %s, Error: %v",
			id, err)

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

func (r *categoryRepository) Create(category *models.Category) error {
	err := r.db.Create(category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperrors.ErrDuplicateEntry
		}

		log.Printf("ERROR: Failed to update category in database. Name: %s, Error: %v",
			category.Name, err)

		return apperrors.ErrDatabaseError
	}
	return nil
}

func (r *categoryRepository) Update(id string, category *models.Category) error {
	category.ID = id
	result := r.db.Model(&models.Category{}).Where("id = ?", id).Updates(category)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return apperrors.ErrDuplicateEntry
		}

		log.Printf("ERROR: Failed to update category in database. Name: %s, Error: %v",
			category.Name, result.Error)

		return apperrors.ErrDatabaseError
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrDuplicateEntry
	}

	return nil
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}
