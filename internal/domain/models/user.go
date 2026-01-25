package models

import (
	"time"

	"github.com/RubenRodrigo/go-tiny-store/internal/infrastructure/db"
)

type User struct {
	db.Base
	Email               string               `json:"email" gorm:"unique;not null"`
	Username            string               `json:"username" gorm:"not null"`
	Password            string               `json:"-" gorm:"not null"` // Password is never returned in JSON
	FirstName           string               `json:"first_name"`
	LastName            string               `json:"last_name"`
	RefreshTokens       []RefreshToken       `json:"refresh_tokens" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PasswordResetTokens []PasswordResetToken `json:"password_reset_tokens" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type RefreshToken struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Token     string    `json:"token" gorm:"unique;not null"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;index"`
}
