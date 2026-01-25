package repository

import (
	"errors"
	"log"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/models"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"gorm.io/gorm"
)

type categoryRepo struct {
	db *gorm.DB
}

type Category interface {
	ListCategory() ([]*models.Category, error)
	GetCategoryById(id string) (*models.Category, error)
	CreateCategory(category *models.Category) error
	UpdateCategory(id string, category *models.Category) error
	DeleteCategory(id string) error
}

func NewCategoryRepository(db *gorm.DB) Category {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) DeleteCategory(id string) error {
	result := r.db.Delete(&models.Category{}, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound
		}
		log.Printf("ERROR: Failed to delete category in database. ID: %s, Error: %v",
			id, result.Error)
		return apperrors.ErrDatabaseError
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *categoryRepo) GetCategoryById(id string) (*models.Category, error) {
	var category models.Category
	if err := r.db.Where("id = ?", id).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		log.Printf("ERROR: Failed to read category in database. ID: %s, Error: %v",
			id, err)
		return nil, apperrors.ErrDatabaseError
	}
	return &category, nil
}

func (r *categoryRepo) ListCategory() ([]*models.Category, error) {
	var categories []*models.Category
	if err := r.db.Find(&categories).Error; err != nil {
		log.Printf("ERROR: Failed to read categories in database. Error: %v", err)
		return nil, apperrors.ErrDatabaseError
	}
	return categories, nil
}

func (r *categoryRepo) CreateCategory(category *models.Category) error {
	err := r.db.Create(category).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperrors.ErrDuplicateEntry
		}
		log.Printf("ERROR: Failed to create category in database. Name: %s, Error: %v",
			category.Name, err)
		return apperrors.ErrDatabaseError
	}
	return nil
}

func (r *categoryRepo) UpdateCategory(id string, category *models.Category) error {
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
		return apperrors.ErrNotFound
	}

	return nil
}
