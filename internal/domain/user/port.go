package user

import "time"

type Service interface {
	Create(input CreateUserInput) (*User, error)
	GetById(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	ListUsers() ([]*User, error)
	IssueRefreshToken(time time.Time, userId, tokenString string) (*RefreshToken, error)
	RevokeToken(tokenString string) error
}
