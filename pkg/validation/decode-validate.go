package validation

import (
	"encoding/json"
	"net/http"

	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
)

func DecodeAndValidate(r *http.Request, v interface{}) error {
	// Decode JSON
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return apperrors.ErrRequestInvalidBody
	}
	// Validate
	if validationErrors := Validate(v); len(validationErrors.Errors) > 0 {
		return validationErrors
	}

	return nil
}
