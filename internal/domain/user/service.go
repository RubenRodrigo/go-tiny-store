package user

import (
	"fmt"
	"time"
)

type service struct {
	repo Repository
}

func NewService(userRepo Repository) Service {
	return &service{
		repo: userRepo,
	}
}

func (s *service) Create(input CreateUserInput) (*User, error) {
	user := &User{
		Email:     input.Email,
		Username:  input.Username,
		Password:  input.PasswordHash,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	err := s.repo.CreateUser(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetById(id string) (*User, error) {
	user, err := s.repo.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetByEmail(email string) (*User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) ListUsers() ([]*User, error) {
	users, err := s.repo.ListUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *service) IssueRefreshToken(time time.Time, userID, tokenString string) (*RefreshToken, error) {
	token := &RefreshToken{
		ExpiresAt: time,
		Token:     tokenString,
		UserID:    userID,
	}

	fmt.Println(token)

	err := s.repo.SaveToken(token)
	if err != nil {
		return nil, err
	}

	return token, err
}

func (s *service) RevokeToken(tokenString string) error {
	return s.repo.DeleteToken(tokenString)
}
