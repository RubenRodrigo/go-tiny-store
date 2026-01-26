package gorm

import (
	"errors"

	"github.com/RubenRodrigo/go-tiny-store/internal/domain/user"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	"gorm.io/gorm"
)

type refreshTokenRepository struct {
	db *gorm.DB
}

// NewRefreshTokenRepository creates a new GORM implementation of user.RefreshTokenRepository
func NewRefreshTokenRepository(db *gorm.DB) user.RefreshTokenRepository {
	return &refreshTokenRepository{db: db}
}

func (r *refreshTokenRepository) SaveToken(rt *user.RefreshToken) error {
	model := toRefreshTokenModel(rt)
	if err := r.db.Create(model).Error; err != nil {
		return apperrors.ErrDatabaseError
	}
	rt.ID = model.ID
	rt.CreatedAt = model.CreatedAt
	rt.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *refreshTokenRepository) GetRefreshToken(token string) (*user.RefreshToken, error) {
	var model RefreshTokenModel
	if err := r.db.Where("token = ?", token).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, apperrors.ErrDatabaseError
	}
	return toRefreshTokenDomain(&model), nil
}

func (r *refreshTokenRepository) DeleteToken(token string) error {
	result := r.db.Where("token = ?", token).Delete(&RefreshTokenModel{})
	if result.Error != nil {
		return apperrors.ErrDatabaseError
	}
	return nil
}

func (r *refreshTokenRepository) DeleteTokensByUserID(userID string) error {
	result := r.db.Where("user_id = ?", userID).Delete(&RefreshTokenModel{})
	if result.Error != nil {
		return apperrors.ErrDatabaseError
	}
	return nil
}

// Mapping functions

func toRefreshTokenModel(rt *user.RefreshToken) *RefreshTokenModel {
	return &RefreshTokenModel{
		ID:        rt.ID,
		Token:     rt.Token,
		ExpiresAt: rt.ExpiresAt,
		UserID:    rt.UserID,
		CreatedAt: rt.CreatedAt,
		UpdatedAt: rt.UpdatedAt,
	}
}

func toRefreshTokenDomain(m *RefreshTokenModel) *user.RefreshToken {
	return &user.RefreshToken{
		ID:        m.ID,
		Token:     m.Token,
		ExpiresAt: m.ExpiresAt,
		UserID:    m.UserID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
