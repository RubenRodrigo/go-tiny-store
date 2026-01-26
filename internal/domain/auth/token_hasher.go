package auth

// TokenHasher defines the interface for hashing tokens (like password reset tokens)
type TokenHasher interface {
	// GenerateToken generates a raw token and its hash
	GenerateToken() (raw string, hash string, err error)

	// HashToken hashes a raw token
	HashToken(raw string) string
}
