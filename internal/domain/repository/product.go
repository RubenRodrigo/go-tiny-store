package repository

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/models"
	"gorm.io/gorm"
)

type productRepo struct {
	db *gorm.DB
}

type Product interface {
	ListProducts() ([]*models.Product, error)
}

func NewProductRepository(db *gorm.DB) Product {
	return &productRepo{db: db}
}

// List implements Repository.
func (r *productRepo) ListProducts() ([]*models.Product, error) {
	panic("unimplemented")
}
