package models

import (
	"time"
)

type PasswordResetToken struct {
	ID        string     `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	TokenHash string     `json:"token_hash" gorm:"unique;not null"`
	ExpiresAt time.Time  `json:"expires_at"`
	UsedAt    *time.Time `json:"used_at"`
	UserID    string     `json:"user_id" gorm:"type:uuid;not null;index"`
}
