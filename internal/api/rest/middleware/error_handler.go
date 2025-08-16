// middleware/error_handler.go
package middleware

import (
	"errors"
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/api/rest/httputil"
	"github.com/RubenRodrigo/go-tiny-store/internal/apperrors"
)

// ErrorHandlerFunc is a type that handles errors and returns HTTP responses
type ErrorHandlerFunc func(r *http.Request, err error) (int, interface{})

// DefaultErrorHandler maps common service errors to HTTP responses
func DefaultErrorHandler(r *http.Request, err error) (int, interface{}) {
	// Check if it's a validation error first
	var validationErrors *apperrors.ValidationErrors
	if errors.As(err, &validationErrors) {
		return http.StatusBadRequest, *validationErrors
	}

	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		return http.StatusNotFound, map[string]string{"error": "Resource not found"}

	case errors.Is(err, apperrors.ErrAuthInvalidCredentials):
		return http.StatusBadRequest, map[string]string{"error": "Invalid Credentials"}

	case errors.Is(err, apperrors.ErrUserEmailExists):
		return http.StatusBadRequest, map[string]string{"error": "This user already exists"}

	case errors.Is(err, apperrors.ErrDatabaseError):
		return http.StatusInternalServerError, map[string]string{"error": "Internal server error"}

	case errors.Is(err, apperrors.ErrRequestInvalidBody):
		return http.StatusInternalServerError, map[string]string{"error": "Invalid request body"}

	// Add other error mappings as needed
	default:
		return http.StatusInternalServerError, map[string]string{"error": "Unexpected error"}
	}
}

// AuthErrorHandler handles authentication-specific errors
func AuthErrorHandler(r *http.Request, err error) (int, interface{}) {
	switch {
	case errors.Is(err, apperrors.ErrAuthMissingToken):
		return http.StatusUnauthorized, map[string]string{"error": "Authorization token required"}
	case errors.Is(err, apperrors.ErrAuthInvalidTokenFormat):
		return http.StatusUnauthorized, map[string]string{"error": "Authorization header format must be 'Bearer {token}'"}
	case errors.Is(err, apperrors.ErrAuthTokenInvalid):
		return http.StatusUnauthorized, map[string]string{"error": "Invalid or malformed token"}
	case errors.Is(err, apperrors.ErrAuthTokenExpired):
		return http.StatusUnauthorized, map[string]string{"error": "Token has expired"}
	case errors.Is(err, apperrors.ErrInsufficientPermissions):
		return http.StatusForbidden, map[string]string{"error": "Insufficient permissions"}
	default:
		// Fall back to default error handler for non-auth errors
		return DefaultErrorHandler(r, err)
	}
}

// WithErrorHandling wraps a handler function that returns an error
func WithErrorHandling(handler func(w http.ResponseWriter, r *http.Request) error, errorHandler ErrorHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			status, message := errorHandler(r, err)
			httputil.RespondWithError(w, status, message)
		}
	}
}
