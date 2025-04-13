// middleware/error_handler.go
package middleware

import (
	"errors"
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/httputil"
	"github.com/RubenRodrigo/go-tiny-store/internal/apperrors"
)

// ErrorHandlerFunc is a type that handles errors and returns HTTP responses
type ErrorHandlerFunc func(err error) (int, string)

// DefaultErrorHandler maps common service errors to HTTP responses
func DefaultErrorHandler(err error) (int, string) {
	switch {
	case errors.Is(err, apperrors.ErrUserNotFound):
		return http.StatusNotFound, "Resource not found"

	case errors.Is(err, apperrors.ErrAuthInvalidCredentials):
		return http.StatusBadRequest, "Invalid Credentials"

	case errors.Is(err, apperrors.ErrAuthEmailExists):
		return http.StatusBadRequest, "This user already exists"

	case errors.Is(err, apperrors.ErrDatabaseError):
		return http.StatusInternalServerError, "Internal server error"

	// Add other error mappings as needed
	default:
		return http.StatusInternalServerError, "Unexpected error"
	}
}

// WithErrorHandling wraps a handler function that returns an error
func WithErrorHandling(handler func(w http.ResponseWriter, r *http.Request) error, errorHandler ErrorHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			status, message := errorHandler(err)
			httputil.RespondWithError(w, status, message)
		}
	}
}
