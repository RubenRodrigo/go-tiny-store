package models

import (
	"time"

	"github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/db"
)

type User struct {
	db.Base
	Email        string         `json:"email" gorm:"unique;not null"`
	Username     string         `json:"username" gorm:"not null"`
	Password     string         `json:"-" gorm:"not null"` // Password is never returned in JSON
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	RefreshToken []RefreshToken `json:"refresh_tokens" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type RefreshToken struct {
	db.Base
	Token     string    `json:"string" gorm:"unique;not null"`
	Revoked   bool      `json:"revoked" gorm:"default:false"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;index"`
}
