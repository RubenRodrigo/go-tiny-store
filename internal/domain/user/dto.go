package user

type CreateUserInput struct {
	Email        string
	Username     string
	FirstName    string
	LastName     string
	PasswordHash string
}
