package gorm

import (
	"errors"
	"log"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/category"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"gorm.io/gorm"
)

type categoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new GORM implementation of category.Repository
func NewCategoryRepository(db *gorm.DB) category.Repository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) ListCategories() ([]*category.Category, error) {
	var models []*CategoryModel
	if err := r.db.Find(&models).Error; err != nil {
		log.Printf("ERROR: Failed to read categories in database. Error: %v", err)
		return nil, apperrors.ErrDatabaseError
	}

	categories := make([]*category.Category, len(models))
	for i, model := range models {
		categories[i] = toCategoryDomain(model)
	}
	return categories, nil
}

func (r *categoryRepository) GetCategoryByID(id string) (*category.Category, error) {
	var model CategoryModel
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		log.Printf("ERROR: Failed to read category in database. ID: %s, Error: %v", id, err)
		return nil, apperrors.ErrDatabaseError
	}
	return toCategoryDomain(&model), nil
}

func (r *categoryRepository) CreateCategory(c *category.Category) error {
	model := toCategoryModel(c)
	err := r.db.Create(model).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperrors.ErrDuplicateEntry
		}
		log.Printf("ERROR: Failed to create category in database. Name: %s, Error: %v", c.Name, err)
		return apperrors.ErrDatabaseError
	}
	c.ID = model.ID
	c.CreatedAt = model.CreatedAt
	c.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *categoryRepository) UpdateCategory(id string, c *category.Category) error {
	c.ID = id
	model := toCategoryModel(c)
	result := r.db.Model(&CategoryModel{}).Where("id = ?", id).Updates(model)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) {
			return apperrors.ErrDuplicateEntry
		}
		log.Printf("ERROR: Failed to update category in database. Name: %s, Error: %v", c.Name, result.Error)
		return apperrors.ErrDatabaseError
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *categoryRepository) DeleteCategory(id string) error {
	result := r.db.Delete(&CategoryModel{}, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound
		}
		log.Printf("ERROR: Failed to delete category in database. ID: %s, Error: %v", id, result.Error)
		return apperrors.ErrDatabaseError
	}

	if result.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

// Mapping functions

func toCategoryModel(c *category.Category) *CategoryModel {
	return &CategoryModel{
		Base: Base{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		},
		Name: c.Name,
	}
}

func toCategoryDomain(m *CategoryModel) *category.Category {
	return &category.Category{
		ID:        m.ID,
		Name:      m.Name,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
