package auth

import "time"

// TokenClaims represents the claims contained in a token (domain value object)
type TokenClaims struct {
	UserID    string
	Email     string
	Username  string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// GeneratedToken represents a generated token with its claims
type GeneratedToken struct {
	Token     string
	UserID    string
	Email     string
	Username  string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// TokenService defines the interface for token operations
type TokenService interface {
	GenerateAccessToken(userID, email, username string) (*GeneratedToken, error)
	GenerateRefreshToken(userID, email, username string) (*GeneratedToken, error)
	ValidateToken(token string) (*TokenClaims, error)
}
