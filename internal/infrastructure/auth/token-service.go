package auth

import "time"

// TokenService defines the interface for token operations
// This makes it easy to swap JWT for another implementation (Paseto, etc.)
type TokenService interface {
	GenerateAccessToken(userID, email, username string) (*TokenGeneratedClaims, error)
	GenerateRefreshToken(userID, email, username string) (*TokenGeneratedClaims, error)
	ValidateToken(token string) (*TokenClaims, error)
}

// TokenClaims represents the parsed JWT claims
type TokenClaims struct {
	UserID    string
	Email     string
	Username  string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

type TokenGeneratedClaims struct {
	UserID    string
	Email     string
	Username  string
	IssuedAt  time.Time
	ExpiresAt time.Time
	Token     string
}
