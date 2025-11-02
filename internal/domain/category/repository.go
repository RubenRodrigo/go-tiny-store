package category

import (
	"errors"
	"log"

	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	List() ([]*Category, error)
	GetById(id string) (*Category, error)
	Create(category *Category) error
	Update(id string, category *Category) error
	Delete(id string) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Delete implements CategoryRepository.
func (r *repository) Delete(id string) error {
	err := r.db.Delete(&Category{}, id).Error

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
func (r *repository) GetById(id string) (*Category, error) {
	var category Category
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
func (r *repository) List() ([]*Category, error) {
	var categories []*Category
	if err := r.db.Find(&categories).Error; err != nil {
		log.Printf("ERROR: Failed to read categories in database. Error: %s", err)
		return nil, apperrors.ErrDatabaseError
	}

	return categories, nil
}

func (r *repository) Create(category *Category) error {
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

func (r *repository) Update(id string, category *Category) error {
	category.ID = id
	result := r.db.Model(&Category{}).Where("id = ?", id).Updates(category)

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
