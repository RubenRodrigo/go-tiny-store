package repository

import (
	"errors"
	"strings"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/models"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"gorm.io/gorm"
)

type User interface {
	CreateUser(user *models.User) error
	GetUserById(id string) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	ListUsers() ([]*models.User, error)
	SaveToken(rt *models.RefreshToken) error
	DeleteToken(token string) error
	UpdateUser(user *models.User) error
}

type userRepo struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) User {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *models.User) error {
	if err := r.db.Create(user).Error; err != nil {
		// Handle duplicate key constraint
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			strings.Contains(err.Error(), "duplicate") ||
			strings.Contains(err.Error(), "unique constraint") {
			if strings.Contains(err.Error(), "email") {
				return apperrors.ErrUserEmailExists
			}
			return apperrors.ErrDuplicateEntry
		}
		return apperrors.ErrDatabaseError
	}
	return nil
}

func (r *userRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrDatabaseError
	}
	return &user, nil
}

func (r *userRepo) GetUserById(id string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrDatabaseError
	}
	return &user, nil
}

func (r *userRepo) ListUsers() ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, apperrors.ErrDatabaseError
	}
	return users, nil
}

func (r *userRepo) UpdateUser(user *models.User) error {
	if err := r.db.Save(user).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) ||
			strings.Contains(err.Error(), "duplicate") {
			if strings.Contains(err.Error(), "email") {
				return apperrors.ErrUserEmailExists
			}
			return apperrors.ErrDuplicateEntry
		}
		return apperrors.ErrDatabaseError
	}
	return nil
}

func (r *userRepo) SaveToken(rt *models.RefreshToken) error {
	if err := r.db.Create(rt).Error; err != nil {
		return apperrors.ErrDatabaseError
	}
	return nil
}

func (r *userRepo) DeleteToken(token string) error {
	// Use Delete with Where to avoid "WHERE conditions required" error
	result := r.db.Where("token = ?", token).Delete(&models.RefreshToken{})
	if result.Error != nil {
		return apperrors.ErrDatabaseError
	}
	// Note: If no rows affected, that's OK (token already deleted or never existed)
	return nil
}
