package user

// Repository defines the interface for user persistence operations
type Repository interface {
	CreateUser(user *User) error
	GetUserByID(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	ListUsers() ([]*User, error)
	UpdateUser(user *User) error
	UpdateUserPassword(userID, newHashedPassword string) error
}

// RefreshTokenRepository defines the interface for refresh token operations
type RefreshTokenRepository interface {
	SaveToken(rt *RefreshToken) error
	GetRefreshToken(token string) (*RefreshToken, error)
	DeleteToken(token string) error
	DeleteTokensByUserID(userID string) error
}

// PasswordResetTokenRepository defines the interface for password reset token operations
type PasswordResetTokenRepository interface {
	CreateToken(token *PasswordResetToken) error
	GetTokenByHash(tokenHash string) (*PasswordResetToken, error)
	MarkTokenAsUsed(tokenHash string) error
	DeleteActiveResetTokens(userID string) error
}
