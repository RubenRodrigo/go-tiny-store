package apperrors

import "errors"

// Common Errors
var (
	ErrDatabaseError  = errors.New("database error")
	ErrDuplicateEntry = errors.New("duplicated entry error")
)

// User related errors
var (
	ErrUserNotFound    = errors.New("user not found")
	ErrUserEmailExists = errors.New("user with this email already exists")
)

// Authentication errors
var (
	ErrAuthInvalidCredentials  = errors.New("invalid credentials")
	ErrAuthUnauthorized        = errors.New("unauthorized access")
	ErrAuthTokenGenerated      = errors.New("failed to generate token")
	ErrAuthTokenInvalid        = errors.New("invalid auth token")
	ErrAuthMissingToken        = errors.New("missing authorization token")
	ErrAuthInvalidTokenFormat  = errors.New("invalid token format")
	ErrAuthTokenExpired        = errors.New("auth token expired")
	ErrInsufficientPermissions = errors.New("insufficient permissions")
)

// Request errors
var (
	ErrRequestInvalidBody = errors.New("invalid request body")
)
