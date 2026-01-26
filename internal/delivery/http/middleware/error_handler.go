package middleware

import (
	"errors"
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
	httputil "github.com/RubenRodrigo/go-tiny-store/pkg/httputils"
)

// HandleError maps AppError and other known errors to HTTP responses
func mapError(r *http.Request, err error) (int, interface{}) {
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

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	status, body := mapError(r, err)
	httputil.RespondWithError(w, status, body)
}
