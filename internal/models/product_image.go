package models

import (
	"time"

	"gorm.io/gorm"
)

type ProductImage struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ProductID string         `json:"product_id" gorm:"not null;index"`
	URL       string         `json:"url" gorm:"not null;size:500"`
	AltText   string         `json:"alt_text" gorm:"size:255"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
