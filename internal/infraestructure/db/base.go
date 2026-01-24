package db

import (
	"time"

	"gorm.io/gorm"
)

type Base struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
