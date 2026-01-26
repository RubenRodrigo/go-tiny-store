package authapp

// SignUpDTO represents the input for user registration
type SignUpDTO struct {
	Email     string
	Username  string
	Password  string
	FirstName string
	LastName  string
}

// SignInDTO represents the input for user login
type SignInDTO struct {
	Email    string
	Password string
}

// AuthUserDTO represents the authenticated user response
type AuthUserDTO struct {
	ID           string `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
