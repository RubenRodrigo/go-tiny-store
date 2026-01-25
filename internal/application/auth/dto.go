package auth

// UserAdapter DTOS
type AuthUserDTO struct {
	ID           string
	Email        string
	Username     string
	FirstName    string
	LastName     string
	PasswordHash string
}
