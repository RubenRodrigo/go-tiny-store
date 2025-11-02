package token

import (
	"time"

	"github.com/RubenRodrigo/go-tiny-store/internal/platform/db"
)

type RefreshToken struct {
	db.Base
	Token     string    `json:"string" gorm:"unique;not null"`
	Revoked   bool      `json:"revoked" gorm:"default:false"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    string    `json:"user_id" gorm:"type:uuid;not null;index"`
}
