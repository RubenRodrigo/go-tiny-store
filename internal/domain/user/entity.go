package user

import "time"

// User represents a user in the system (pure domain entity)
type User struct {
	ID        string
	Email     string
	Username  string
	Password  string
	FirstName string
	LastName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Role represents a user role
type Role struct {
	ID          string
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// RefreshToken represents a refresh token for authentication
type RefreshToken struct {
	ID        string
	Token     string
	ExpiresAt time.Time
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// PasswordResetToken represents a token for password reset
type PasswordResetToken struct {
	ID        string
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	UserID    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
