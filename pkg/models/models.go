package models

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/category"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/product"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/product_image"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/user"
)

// AllModels returns all models for schema generation.
func AllModels() []interface{} {
	return []interface{}{
		&user.User{},
		&category.Category{},
		&product.Product{},
		&product_image.ProductImage{},
	}
}
