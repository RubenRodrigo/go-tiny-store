package httputil

import (
	"context"
	"strconv"

	"github.com/RubenRodrigo/go-tiny-store/internal/apperrors"
)

// GetUserIDFromContext extracts user ID from request context
func GetUserIDFromContext(ctx context.Context) (uint, error) {
	userID := ctx.Value("userID")
	if userID == nil {
		return 0, apperrors.ErrAuthUnauthorized
	}

	// Handle different types that might be stored
	switch v := userID.(type) {
	case string:
		id, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return 0, apperrors.ErrAuthTokenInvalid
		}
		return uint(id), nil
	case uint:
		return v, nil
	case float64: // JWT claims are often float64
		return uint(v), nil
	default:
		return 0, apperrors.ErrAuthTokenInvalid
	}
}

// GetEmailFromContext extracts email from request context
func GetEmailFromContext(ctx context.Context) (string, error) {
	email := ctx.Value("email")
	if email == nil {
		return "", apperrors.ErrAuthTokenInvalid
	}

	emailStr, ok := email.(string)
	if !ok {
		return "", apperrors.ErrAuthTokenInvalid
	}

	return emailStr, nil
}

// GetUsernameFromContext extracts username from request context
func GetUsernameFromContext(ctx context.Context) (string, error) {
	username := ctx.Value("username")
	if username == nil {
		return "", apperrors.ErrAuthTokenInvalid
	}

	usernameStr, ok := username.(string)
	if !ok {
		return "", apperrors.ErrAuthTokenInvalid
	}

	return usernameStr, nil
}

// AuthUser represents the authenticated user info
type AuthUser struct {
	ID       uint   `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

// GetAuthUserFromContext gets all auth info as a struct
func GetAuthUserFromContext(ctx context.Context) (*AuthUser, error) {
	userID, err := GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	email, _ := GetEmailFromContext(ctx)       // Optional
	username, _ := GetUsernameFromContext(ctx) // Optional

	return &AuthUser{
		ID:       userID,
		Email:    email,
		Username: username,
	}, nil
}
