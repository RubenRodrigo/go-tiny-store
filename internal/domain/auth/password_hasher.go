package auth

// PasswordHasher defines the interface for password hashing operations
type PasswordHasher interface {
	// HashPassword generates a hash from a plain password
	HashPassword(password string) (string, error)

	// ComparePassword compares a plain password with a hashed password
	ComparePassword(hashedPassword, password string) error
}
