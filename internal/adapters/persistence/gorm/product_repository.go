package gorm

import (
	"errors"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/product"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"github.com/RubenRodrigo/go-tiny-store/pkg/pagination"
	"gorm.io/gorm"
)

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new GORM implementation of product.Repository
func NewProductRepository(db *gorm.DB) product.Repository {
	return &productRepository{db: db}
}

func (r *productRepository) ListProducts(params pagination.Params, filters product.Filters) ([]*product.Product, int64, error) {
	var models []*ProductModel
	var totalCount int64

	// Build query with filters
	query := r.db.Model(&ProductModel{})
	query = applyProductFilters(query, filters)

	// Count total records (with filters applied)
	if err := query.Count(&totalCount).Error; err != nil {
		return nil, 0, apperrors.ErrDatabaseError
	}

	// Apply pagination and fetch results
	err := query.
		Preload("Images").
		Offset(params.Offset()).
		Limit(params.Limit()).
		Find(&models).Error

	if err != nil {
		return nil, 0, apperrors.ErrDatabaseError
	}

	products := make([]*product.Product, len(models))
	for i, model := range models {
		products[i] = toProductDomain(model)
	}

	return products, totalCount, nil
}

// applyProductFilters applies filters to the GORM query
func applyProductFilters(query *gorm.DB, filters product.Filters) *gorm.DB {
	// Filter by category
	if filters.CategoryID != "" {
		query = query.Where("category_id = ?", filters.CategoryID)
	}

	// Filter by minimum price
	if filters.MinPrice != nil {
		query = query.Where("price >= ?", *filters.MinPrice)
	}

	// Filter by maximum price
	if filters.MaxPrice != nil {
		query = query.Where("price <= ?", *filters.MaxPrice)
	}

	// Filter by disabled status
	if filters.Disabled != nil {
		query = query.Where("disabled = ?", *filters.Disabled)
	}

	return query
}

func (r *productRepository) GetProduct(id string) (*product.Product, error) {
	var model ProductModel
	if err := r.db.Preload("Images").First(&model, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrDatabaseError
	}

	return toProductDomain(&model), nil
}

func (r *productRepository) CreateProduct(p *product.Product) error {
	model := toProductModel(p)
	if err := r.db.Create(model).Error; err != nil {
		return apperrors.ErrDatabaseError
	}
	p.ID = model.ID
	p.CreatedAt = model.CreatedAt
	p.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *productRepository) UpdateProduct(p *product.Product) error {
	model := toProductModel(p)
	if err := r.db.Model(&ProductModel{}).Where("id = ?", p.ID).Updates(model).Error; err != nil {
		return apperrors.ErrDatabaseError
	}
	return nil
}

// Mapping functions

func toProductModel(p *product.Product) *ProductModel {
	images := make([]ProductImageModel, len(p.Images))
	for i, img := range p.Images {
		images[i] = ProductImageModel{
			Base: Base{
				ID:        img.ID,
				CreatedAt: img.CreatedAt,
				UpdatedAt: img.UpdatedAt,
			},
			ProductID: img.ProductID,
			URL:       img.URL,
			AltText:   img.AltText,
		}
	}

	return &ProductModel{
		Base: Base{
			ID:        p.ID,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		},
		Name:       p.Name,
		Price:      p.Price,
		Disabled:   p.Disabled,
		Stock:      p.Stock,
		CategoryID: p.CategoryID,
		Images:     images,
	}
}

func toProductDomain(m *ProductModel) *product.Product {
	images := make([]product.ProductImage, len(m.Images))
	for i, img := range m.Images {
		images[i] = product.ProductImage{
			ID:        img.ID,
			ProductID: img.ProductID,
			URL:       img.URL,
			AltText:   img.AltText,
			CreatedAt: img.CreatedAt,
			UpdatedAt: img.UpdatedAt,
		}
	}

	return &product.Product{
		ID:         m.ID,
		Name:       m.Name,
		Price:      m.Price,
		Disabled:   m.Disabled,
		Stock:      m.Stock,
		CategoryID: m.CategoryID,
		Images:     images,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}
