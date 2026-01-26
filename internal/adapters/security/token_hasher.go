package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/auth"
)

type sha256TokenHasher struct{}

// NewSHA256TokenHasher creates a new SHA256-based token hasher
func NewSHA256TokenHasher() auth.TokenHasher {
	return &sha256TokenHasher{}
}

func (h *sha256TokenHasher) GenerateToken() (raw string, hash string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", "", err
	}

	raw = base64.RawURLEncoding.EncodeToString(b)
	hash = h.HashToken(raw)

	return raw, hash, nil
}

func (h *sha256TokenHasher) HashToken(raw string) string {
	hashBytes := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(hashBytes[:])
}
