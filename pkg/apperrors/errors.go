package apperrors

import "net/http"

// AppError defines a structured application error
type AppError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Status  int         `json:"-"`
	Details interface{} `json:"details,omitempty"`
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// New creates a new AppError
func New(code, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// Common errors
var (
	ErrDatabaseError  = New("DATABASE_ERROR", "A database error occurred", http.StatusInternalServerError)
	ErrDuplicateEntry = New("DUPLICATE_ENTRY", "Duplicated entry", http.StatusConflict)
)

// User-related errors
var (
	ErrUserEmailExists = New("USER_EMAIL_EXISTS", "User with this email already exists", http.StatusBadRequest)
)

// Authentication errors
var (
	ErrAuthInvalidCredentials  = New("INVALID_CREDENTIALS", "Invalid credentials", http.StatusBadRequest)
	ErrAuthUnauthorized        = New("UNAUTHORIZED_ACCESS", "Unauthorized access", http.StatusUnauthorized)
	ErrAuthTokenGenerated      = New("TOKEN_GENERATION_FAILED", "Failed to generate token", http.StatusInternalServerError)
	ErrAuthTokenInvalid        = New("INVALID_TOKEN", "Invalid authentication token", http.StatusUnauthorized)
	ErrAuthMissingToken        = New("MISSING_TOKEN", "Missing authorization token", http.StatusUnauthorized)
	ErrAuthInvalidTokenFormat  = New("INVALID_TOKEN_FORMAT", "Invalid token format", http.StatusUnauthorized)
	ErrAuthInvalidResetToken   = New("INVALID_RESET_TOKEN", "Invalid reset token", http.StatusUnauthorized)
	ErrAuthTokenExpired        = New("TOKEN_EXPIRED", "Authentication token expired", http.StatusUnauthorized)
	ErrInsufficientPermissions = New("INSUFFICIENT_PERMISSIONS", "You do not have permission to perform this action", http.StatusForbidden)
)

// Request errors
var (
	ErrRequestInvalidBody = New("INVALID_REQUEST_BODY", "Invalid request body", http.StatusBadRequest)
)

// Resource errors
var (
	ErrNotFound      = New("NOT_FOUND", "Resource not found", http.StatusNotFound)
	ErrAlreadyExists = New("ALREADY_EXISTS", "Resource already exists", http.StatusConflict)
)
