package repository

import (
	"errors"
	"time"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/models"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"gorm.io/gorm"
)

type PasswordResetToken interface {
	CreateToken(token *models.PasswordResetToken) error
	GetTokenByHash(tokenHash string) (*models.PasswordResetToken, error)
	MarkTokenAsUsed(token string) error
	DeleteActiveResetTokens(userId string) error
}

type passwordResetTokenRepo struct {
	db *gorm.DB
}

func NewPasswordResetTokenRepository(db *gorm.DB) PasswordResetToken {
	return &passwordResetTokenRepo{db: db}
}

func (r *passwordResetTokenRepo) CreateToken(token *models.PasswordResetToken) error {
	if err := r.db.Create(token).Error; err != nil {
		return apperrors.ErrDatabaseError
	}
	return nil
}

func (r *passwordResetTokenRepo) GetTokenByHash(tokenHash string) (*models.PasswordResetToken, error) {
	var prt models.PasswordResetToken
	now := time.Now()

	if err := r.db.Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", tokenHash, now).First(&prt).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrDatabaseError
	}
	return &prt, nil
}

func (r *passwordResetTokenRepo) MarkTokenAsUsed(token string) error {
	result := r.db.Model(&models.PasswordResetToken{}).
		Where("token_hash = ?", token).
		Update("used_at", time.Now())

	if result.Error != nil {
		return apperrors.ErrDatabaseError
	}
	return nil
}

func (r *passwordResetTokenRepo) DeleteActiveResetTokens(userId string) error {
	result := r.db.Where("user_id = ? AND used_at IS NULL", userId).Delete(&models.PasswordResetToken{})
	if result.Error != nil {
		return apperrors.ErrDatabaseError
	}

	return nil
}
