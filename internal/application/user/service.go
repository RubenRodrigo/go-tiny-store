package user

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/models"
	"github.com/RubenRodrigo/go-tiny-store/internal/domain/repository"
)

type Service struct {
	repo repository.User
}

func NewService(userRepo repository.User) *Service {
	return &Service{
		repo: userRepo,
	}
}

func (s *Service) Create(input CreateUserInput) (*models.User, error) {
	user := &models.User{
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

func (s *Service) GetById(id string) (*models.User, error) {
	user, err := s.repo.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) GetByEmail(email string) (*models.User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) ListUsers() ([]*models.User, error) {
	users, err := s.repo.ListUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}
