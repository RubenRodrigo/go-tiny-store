package models

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/db"
)

type Category struct {
	db.Base
	Name     string    `json:"name" gorm:"not null;size:100"`
	Products []Product `json:"products"`
}
