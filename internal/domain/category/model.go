package category

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/product"
	"github.com/RubenRodrigo/go-tiny-store/internal/platform/db"
)

type Category struct {
	db.Base
	Name     string            `json:"name" gorm:"not null;size:100"`
	Products []product.Product `json:"products"`
}
