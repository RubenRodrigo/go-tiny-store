package auth

import (
	"time"
)

type Service interface {
	RegisterUser(email, username, password, firstName, lastName string) (*AuthUserDTO, error)
	LoginUser(email, password string) (*AuthUserDTO, string, string, error)
	LogOutUser(token string) error
	RefreshToken(token string) (*AuthUserDTO, string, string, error)
}

type UserServicePort interface {
	GetByEmail(email string) (*AuthUserDTO, error)
	GetById(id string) (*AuthUserDTO, error)
	Create(input RegisterUserDto) (*AuthUserDTO, error)
	IssueRefreshToken(time time.Time, userId, token string) error
	RevokeToken(token string) error
}
