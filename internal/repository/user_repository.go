package repository

import (
	"errors"
	"log"
	"strings"

	"github.com/RubenRodrigo/go-tiny-store/internal/apperrors"
	"github.com/RubenRodrigo/go-tiny-store/internal/models"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *models.User) error {
	err := r.db.Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			if strings.Contains(err.Error(), "email") {
				return apperrors.ErrAuthEmailExists
			}

			return apperrors.ErrDuplicateEntry
		}

		log.Printf("ERROR: Failed to create user in database. Email: %s, Error: %v",
			user.Email, err)

		return apperrors.ErrDatabaseError
	}

	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}

		log.Printf("ERROR: Failed to read user in database. Email: %s, Error: %v",
			email, err)

		return nil, apperrors.ErrDatabaseError
	}

	return &user, nil
}

func (r *userRepository) GetUserById(id string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrUserNotFound
		}

		log.Printf("ERROR: Failed to read user in database. Id: %s, Error: %v",
			id, err)

		return nil, apperrors.ErrDatabaseError
	}

	return &user, nil
}

func (r *userRepository) ListUsers() ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Find(&users).Error; err != nil {
		log.Printf("ERROR: Failed to read users in database. Error: %s", err)
		return nil, apperrors.ErrDatabaseError
	}

	return users, nil
}
