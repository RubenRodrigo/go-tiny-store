package service

import (
	"github.com/RubenRodrigo/go-tiny-store/internal/models"
)

type UserService interface {
	GetUserByID(id string) (*models.User, error)
	ListUsers() ([]*models.User, error)
}

type AuthService interface {
	RegisterUser(email, username, password, firstName, lastName string) (*models.User, error)
	LoginUser(email, password string) (*models.User, string, error)
	LogOutUser(token string) error
}

type CategoryService interface {
	List() ([]*models.Category, error)
	Save(name string, id uint) (*models.Category, error)
	Delete(id uint) (*models.Category, error)
}
