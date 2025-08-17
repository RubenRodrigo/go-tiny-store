package models

type Product struct {
	Base
	Name       string         `json:"name" gorm:"not null;size:255"`
	Price      float64        `json:"price" gorm:"not null"`
	Disabled   bool           `json:"disabled" gorm:"default:false"`
	Stock      int            `json:"stock" gorm:"default:0"`
	CategoryID string         `json:"category_id" gorm:"type:uuid"`
	Images     []ProductImage `json:"images" gorm:"foreignKey:ProductID"`
}
