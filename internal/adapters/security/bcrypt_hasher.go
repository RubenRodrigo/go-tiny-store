package security

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/auth"
	"golang.org/x/crypto/bcrypt"
)

type bcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a new bcrypt-based password hasher
func NewBcryptHasher() auth.PasswordHasher {
	return &bcryptHasher{
		cost: bcrypt.DefaultCost,
	}
}

func (h *bcryptHasher) HashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func (h *bcryptHasher) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
