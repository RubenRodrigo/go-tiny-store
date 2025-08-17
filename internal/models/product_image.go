package models

type ProductImage struct {
	Base
	ProductID string `json:"product_id" gorm:"type:uuid;not null;index"`
	URL       string `json:"url" gorm:"not null;size:500"`
	AltText   string `json:"alt_text" gorm:"size:255"`
}
