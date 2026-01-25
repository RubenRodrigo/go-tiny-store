package auth

type AuthUserDTO struct {
	ID           string
	Email        string
	Username     string
	FirstName    string
	LastName     string
	AccessToken  string
	RefreshToken string
}
