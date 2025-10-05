package user

import (
	"errors"
	"log"
	"strings"

	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	CreateUser(user *User) error
	GetUserById(id string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	ListUsers() ([]*User, error)
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(user *User) error {
	err := r.db.Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			if strings.Contains(err.Error(), "email") {
				return apperrors.ErrUserEmailExists
			}

			return apperrors.ErrDuplicateEntry
		}

		log.Printf("ERROR: Failed to create user in database. Email: %s, Error: %v",
			user.Email, err)

		return apperrors.ErrDatabaseError
	}

	return nil
}

func (r *repository) GetUserByEmail(email string) (*User, error) {
	var user User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}

		log.Printf("ERROR: Failed to read user in database. Email: %s, Error: %v",
			email, err)

		return nil, apperrors.ErrDatabaseError
	}

	return &user, nil
}

func (r *repository) GetUserById(id string) (*User, error) {
	var user User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}

		log.Printf("ERROR: Failed to read user in database. Id: %s, Error: %v",
			id, err)

		return nil, apperrors.ErrDatabaseError
	}

	return &user, nil
}

func (r *repository) ListUsers() ([]*User, error) {
	var users []*User
	if err := r.db.Find(&users).Error; err != nil {
		log.Printf("ERROR: Failed to read users in database. Error: %s", err)
		return nil, apperrors.ErrDatabaseError
	}

	return users, nil
}
