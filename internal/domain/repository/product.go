package repository

import (
	"errors"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/models"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"github.com/RubenRodrigo/go-tiny-store/pkg/pagination"
	"gorm.io/gorm"
)

type productRepo struct {
	db *gorm.DB
}

type Product interface {
	ListProducts(pagination pagination.Params) ([]*models.Product, int64, error)
	GetProduct(id string) (*models.Product, error)
	CreateProduct(product *models.Product) error
}

func NewProductRepository(db *gorm.DB) Product {
	return &productRepo{db: db}
}

// List implements Repository.
func (r *productRepo) ListProducts(pagination pagination.Params) ([]*models.Product, int64, error) {
	var products []*models.Product
	var totalCount int64

	if err := r.db.Model(&models.Product{}).Count(&totalCount).Error; err != nil {
		return nil, 0, apperrors.ErrDatabaseError
	}

	err := r.db.
		Preload("Images").
		Offset(pagination.Offset()).
		Limit(pagination.Limit()).
		Find(&products).Error

	return products, totalCount, err
}

func (r *productRepo) GetProduct(id string) (*models.Product, error) {
	var product models.Product
	if err := r.db.Preload("Images").First(&product, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrDatabaseError
	}

	return &product, nil
}

func (r *productRepo) CreateProduct(product *models.Product) error {
	if err := r.db.Create(product).Error; err != nil {
		return apperrors.ErrDatabaseError
	}

	return nil
}

func (r *productRepo) UpdateProduct(product *models.Product) error {
	if err := r.db.Model(&models.Product{}).Where("id = ?", product.ID).Updates(product).Error; err != nil {
		return apperrors.ErrDatabaseError
	}

	return nil
}
