package repository

import "github.com/RubenRodrigo/go-tiny-store/internal/models"

type UserRepository interface {
	CreateUser(user *models.User) error
	ListUsers() ([]*models.User, error)
	GetUserById(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
}
