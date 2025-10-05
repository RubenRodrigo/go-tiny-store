package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/RubenRodrigo/go-tiny-store/pkg/apperrors"
)

func Validate(val interface{}) *apperrors.ValidationErrors {
	validationErrors := apperrors.NewValidationError()

	v := reflect.ValueOf(val)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).String()
		jsonTag := field.Tag.Get("json")
		validateTag := field.Tag.Get("validate")

		// If it's empty, skip
		if validateTag == "" {
			continue
		}

		fieldName := jsonTag
		if fieldName == "" {
			fieldName = strings.ToLower(field.Name)
		} else {
			// Handle cases like `json:"email,omitempty"`
			fieldName = strings.Split(fieldName, ",")[0]
		}

		rules := strings.Split(validateTag, ",")

		for _, rule := range rules {
			switch {

			case rule == "required":
				if strings.TrimSpace(value) == "" {
					validationErrors.Add(fieldName, fmt.Sprintf("%s is required", fieldName))
				}

			case value != "" && strings.HasPrefix(rule, "min="):
				min, _ := strconv.Atoi(strings.TrimPrefix(rule, "min="))
				if len(value) < min {
					validationErrors.Add(fieldName, fmt.Sprintf("%s must be at least %d characters", fieldName, min))
				}

			case value != "" && strings.HasPrefix(rule, "max="):
				max, _ := strconv.Atoi(strings.TrimPrefix(rule, "max="))
				if len(value) > max {
					validationErrors.Add(fieldName, fmt.Sprintf("%s cannot exceed %d characters", fieldName, max))
				}

			case value != "" && rule == "email":
				emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

				if !emailRegex.MatchString(value) {
					validationErrors.Add(fieldName, "invalid email format")
				}

			}
		}
	}

	return validationErrors
}
