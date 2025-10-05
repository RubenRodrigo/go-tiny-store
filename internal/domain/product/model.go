package product

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/product_image"
	"github.com/RubenRodrigo/go-tiny-store/internal/platform/db"
)

type Product struct {
	db.Base
	Name       string                       `json:"name" gorm:"not null;size:255"`
	Price      float64                      `json:"price" gorm:"not null"`
	Disabled   bool                         `json:"disabled" gorm:"default:false"`
	Stock      int                          `json:"stock" gorm:"default:0"`
	CategoryID string                       `json:"category_id" gorm:"type:uuid"`
	Images     []product_image.ProductImage `json:"images" gorm:"foreignKey:ProductID"`
}
