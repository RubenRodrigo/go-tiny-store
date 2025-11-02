package user

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/token"
	"github.com/RubenRodrigo/go-tiny-store/internal/platform/db"
)

type User struct {
	db.Base
	Email        string               `json:"email" gorm:"unique;not null"`
	Username     string               `json:"username" gorm:"not null"`
	Password     string               `json:"-" gorm:"not null"` // Password is never returned in JSON
	FirstName    string               `json:"first_name"`
	LastName     string               `json:"last_name"`
	RefreshToken []token.RefreshToken `json:"refresh_tokens" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
