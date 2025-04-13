package apperrors

import "errors"

// Common Errors
var (
	ErrDatabaseError  = errors.New("database error")
	ErrDuplicateEntry = errors.New("duplicated entry error")
)

// User related errors
var (
	ErrUserNotFound = errors.New("user not found")
)

// Authentication errors
var (
	ErrAuthEmailExists        = errors.New("user with this email already exists")
	ErrAuthInvalidCredentials = errors.New("invalid credentials")
	ErrAuthUnauthorized       = errors.New("unauthorized access")
	ErrAuthTokenExpired       = errors.New("token has expired")
	ErrAuthTokenGenerated     = errors.New("failed to generate token")
	ErrAuthRequiredFields     = errors.New("email and password are required")
)

// Request errors
var (
	ErrRequestInvalidBody = errors.New("invalid request body")
)
