package models

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/db"
)

type Product struct {
	db.Base
	Name       string         `json:"name" gorm:"not null;size:255"`
	Price      float64        `json:"price" gorm:"not null"`
	Disabled   bool           `json:"disabled" gorm:"default:false"`
	Stock      int            `json:"stock" gorm:"default:0"`
	CategoryID string         `json:"category_id" gorm:"type:uuid"`
	Images     []ProductImage `json:"images" gorm:"foreignKey:ProductID"`
}

type ProductImage struct {
	db.Base
	ProductID string `json:"product_id" gorm:"type:uuid;not null;index"`
	URL       string `json:"url" gorm:"not null;size:500"`
	AltText   string `json:"alt_text" gorm:"size:255"`
}
