package service

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/models"
	"github.com/RubenRodrigo/go-tiny-store/internal/repository"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByID(id string) (*models.User, error) {
	user, err := s.userRepo.GetUserById(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) ListUsers() ([]*models.User, error) {
	users, err := s.userRepo.ListUsers()
	if err != nil {
		return nil, err
	}

	return users, nil
}
