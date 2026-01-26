package gorm

import (
	"errors"
	"time"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/user"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"gorm.io/gorm"
)

type passwordResetTokenRepository struct {
	db *gorm.DB
}

// NewPasswordResetTokenRepository creates a new GORM implementation of user.PasswordResetTokenRepository
func NewPasswordResetTokenRepository(db *gorm.DB) user.PasswordResetTokenRepository {
	return &passwordResetTokenRepository{db: db}
}

func (r *passwordResetTokenRepository) CreateToken(token *user.PasswordResetToken) error {
	model := toPasswordResetTokenModel(token)
	if err := r.db.Create(model).Error; err != nil {
		return apperrors.ErrDatabaseError
	}
	token.ID = model.ID
	token.CreatedAt = model.CreatedAt
	token.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *passwordResetTokenRepository) GetTokenByHash(tokenHash string) (*user.PasswordResetToken, error) {
	var model PasswordResetTokenModel
	now := time.Now()

	if err := r.db.Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", tokenHash, now).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrDatabaseError
	}
	return toPasswordResetTokenDomain(&model), nil
}

func (r *passwordResetTokenRepository) MarkTokenAsUsed(tokenHash string) error {
	result := r.db.Model(&PasswordResetTokenModel{}).
		Where("token_hash = ?", tokenHash).
		Update("used_at", time.Now())

	if result.Error != nil {
		return apperrors.ErrDatabaseError
	}
	return nil
}

func (r *passwordResetTokenRepository) DeleteActiveResetTokens(userID string) error {
	result := r.db.Where("user_id = ? AND used_at IS NULL", userID).Delete(&PasswordResetTokenModel{})
	if result.Error != nil {
		return apperrors.ErrDatabaseError
	}
	return nil
}

// Mapping functions

func toPasswordResetTokenModel(t *user.PasswordResetToken) *PasswordResetTokenModel {
	return &PasswordResetTokenModel{
		ID:        t.ID,
		TokenHash: t.TokenHash,
		ExpiresAt: t.ExpiresAt,
		UsedAt:    t.UsedAt,
		UserID:    t.UserID,
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}

func toPasswordResetTokenDomain(m *PasswordResetTokenModel) *user.PasswordResetToken {
	return &user.PasswordResetToken{
		ID:        m.ID,
		TokenHash: m.TokenHash,
		ExpiresAt: m.ExpiresAt,
		UsedAt:    m.UsedAt,
		UserID:    m.UserID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
