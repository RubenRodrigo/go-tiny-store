package token

import (
	"errors"
	"log"

	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

type Repository interface {
	Create(refreshToken *RefreshToken) error
	DeleteByToken(token string) error
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(refreshToken *RefreshToken) error {
	err := r.db.Create(refreshToken).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return apperrors.ErrDuplicateEntry
		}

		log.Printf("ERROR: Failed to create a refresh token in database. User: %s, Error: %v",
			refreshToken.UserID, err)

		return apperrors.ErrDatabaseError
	}
	return nil
}

func (r *repository) DeleteByToken(token string) error {
	err := r.db.Where("token = ?", token).Delete(&RefreshToken{}).Error

	if err == nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperrors.ErrNotFound
		}

		log.Printf("ERROR: Failed to delete token in database. Error: %v", err)

		return apperrors.ErrDatabaseError
	}

	return nil

}
