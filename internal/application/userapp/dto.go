package userapp

// CreateUserInput represents the input for creating a user
type CreateUserInput struct {
	Email        string
	Username     string
	PasswordHash string
	FirstName    string
	LastName     string
}
