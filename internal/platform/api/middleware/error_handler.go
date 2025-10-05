package middleware

import (
	"errors"
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/internal/platform/api/httputil"
	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
)

// ErrorHandlerFunc defines a handler that maps errors to HTTP responses
type ErrorHandlerFunc func(r *http.Request, err error) (int, interface{})

// DefaultErrorHandler maps AppError and other known errors to HTTP responses
func DefaultErrorHandler(r *http.Request, err error) (int, interface{}) {
	// Handle structured AppError
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		return appErr.Status, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    appErr.Code,
				"message": appErr.Message,
				"details": appErr.Details,
			},
		}
	}

	// Handle validation errors (if any)
	var validationErrors *apperrors.ValidationErrors
	if errors.As(err, &validationErrors) {
		return http.StatusBadRequest, map[string]interface{}{
			"error": map[string]interface{}{
				"code":    "VALIDATION_ERROR",
				"message": "Validation failed",
				"details": *validationErrors,
			},
		}
	}

	// Unknown error
	return http.StatusInternalServerError, map[string]interface{}{
		"error": map[string]interface{}{
			"code":    "INTERNAL_ERROR",
			"message": "Unexpected internal server error",
		},
	}
}

// WithErrorHandling wraps handlers to catch and process errors centrally
func WithErrorHandling(handler func(w http.ResponseWriter, r *http.Request) error, errorHandler ErrorHandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			status, body := errorHandler(r, err)
			httputil.RespondWithError(w, status, body)
		}
	}
}
