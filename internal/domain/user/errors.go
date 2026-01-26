package user

import "github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"

var (
	// ErrEmailAlreadyExists indicates that a user with the given email already exists
	ErrEmailAlreadyExists = apperrors.ErrUserEmailExists

	// ErrNotFound indicates that the requested user was not found
	ErrNotFound = apperrors.ErrNotFound

	// ErrInvalidCredentials indicates that the provided credentials are invalid
	ErrInvalidCredentials = apperrors.ErrAuthInvalidCredentials
)
