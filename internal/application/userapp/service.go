package userapp

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/user"
)

// Service handles user-related use cases
type Service struct {
	userRepo user.Repository
}

// NewService creates a new user application service
func NewService(userRepo user.Repository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

func (s *Service) Create(input CreateUserInput) (*user.User, error) {
	u := &user.User{
		Email:     input.Email,
		Username:  input.Username,
		Password:  input.PasswordHash,
		FirstName: input.FirstName,
		LastName:  input.LastName,
	}

	err := s.userRepo.CreateUser(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (s *Service) GetByID(id string) (*user.User, error) {
	return s.userRepo.GetUserByID(id)
}

func (s *Service) GetByEmail(email string) (*user.User, error) {
	return s.userRepo.GetUserByEmail(email)
}

func (s *Service) ListUsers() ([]*user.User, error) {
	return s.userRepo.ListUsers()
}
