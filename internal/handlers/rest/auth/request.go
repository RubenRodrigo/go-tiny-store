package authhttp

type RegisterUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"min=3,max=30"`
	Password  string `json:"password"  validate:"required,min=8,max=30"`
	FirstName string `json:"first_name" validate:"min=2"`
	LastName  string `json:"last_name" validate:"min=2"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LogOutUserRequest struct {
	Token string `json:"token" validate:"required"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"string" validate:"required"`
}
