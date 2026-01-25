package authhttp

type SignUpRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"min=3,max=30"`
	Password  string `json:"password"  validate:"required,min=8,max=30"`
	FirstName string `json:"first_name" validate:"min=2"`
	LastName  string `json:"last_name" validate:"min=2"`
}

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type SignOutRequest struct {
	Token string `json:"token" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required"`
}

type ResetPasswordRequest struct {
	Token    string `json:"token" validate:"required"`
	Password string `json:"password" validate:"required"`
}
