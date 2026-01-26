package gorm

import (
	"errors"
	"strings"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/user"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new GORM implementation of user.Repository
func NewUserRepository(db *gorm.DB) user.Repository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(u *user.User) error {
	model := toUserModel(u)
	if err := r.db.Create(model).Error; err != nil {
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
	// Update the domain entity with generated ID and timestamps
	u.ID = model.ID
	u.CreatedAt = model.CreatedAt
	u.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *userRepository) GetUserByEmail(email string) (*user.User, error) {
	var model UserModel
	if err := r.db.Where("email = ?", email).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrDatabaseError
	}
	return toUserDomain(&model), nil
}

func (r *userRepository) GetUserByID(id string) (*user.User, error) {
	var model UserModel
	if err := r.db.Where("id = ?", id).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrDatabaseError
	}
	return toUserDomain(&model), nil
}

func (r *userRepository) ListUsers() ([]*user.User, error) {
	var models []*UserModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, apperrors.ErrDatabaseError
	}

	users := make([]*user.User, len(models))
	for i, model := range models {
		users[i] = toUserDomain(model)
	}
	return users, nil
}

func (r *userRepository) UpdateUser(u *user.User) error {
	model := toUserModel(u)
	if err := r.db.Model(&UserModel{}).Where("id = ?", u.ID).Updates(model).Error; err != nil {
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

func (r *userRepository) UpdateUserPassword(userID, newHashedPassword string) error {
	err := r.db.Model(&UserModel{}).Where("id = ?", userID).Update("password", newHashedPassword)
	if err.Error != nil {
		return apperrors.ErrDatabaseError
	}
	return nil
}

// Mapping functions: Domain <-> GORM Model

func toUserModel(u *user.User) *UserModel {
	return &UserModel{
		Base: Base{
			ID:        u.ID,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		},
		Email:     u.Email,
		Username:  u.Username,
		Password:  u.Password,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
}

func toUserDomain(m *UserModel) *user.User {
	return &user.User{
		ID:        m.ID,
		Email:     m.Email,
		Username:  m.Username,
		Password:  m.Password,
		FirstName: m.FirstName,
		LastName:  m.LastName,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
