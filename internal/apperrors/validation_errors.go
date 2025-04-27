package apperrors

import "fmt"

// ValidationError represents a single field validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationErrors holds multiple validation errors
type ValidationErrors struct {
	Errors []ValidationError `json:"errors"`
}

// NewValidationError creates a new ValidationErrors instance
func NewValidationError() *ValidationErrors {
	return &ValidationErrors{
		Errors: []ValidationError{},
	}
}

// Error implements the error interface
func (ve *ValidationErrors) Error() string {
	return fmt.Sprintf("validation failed with %d errors", len(ve.Errors))
}

// Add appends a new validation error
func (ve *ValidationErrors) Add(field, message string) {
	ve.Errors = append(ve.Errors, ValidationError{Field: field, Message: message})
}
