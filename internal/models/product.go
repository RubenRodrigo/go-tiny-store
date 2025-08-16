package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	Name       string         `json:"name" gorm:"not null;size:255"`
	Price      float64        `json:"price" gorm:"not null"`
	Disabled   bool           `json:"disabled" gorm:"default:false"`
	Images     []ProductImage `json:"images" gorm:"foreignKey:ProductID"`
	Stock      int            `json:"stock" gorm:"default:0"`
	CategoryID uint           `json:"category_id"`
	Category   Category       `json:"category" gorm:"foreignKey:CategoryID"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
}
