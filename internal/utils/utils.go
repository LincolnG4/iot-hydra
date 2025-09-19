package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

// FormatValidationErrors takes a validation error and returns a formatted,
// human-readable string with details about each validation failure.
// If the error is not a validator.ValidationErrors, it returns the original error message.
func FormatValidationErrors(err error) string {
	// If there's no error, return an empty string
	if err == nil {
		return ""
	}

	// Use errors.As to check if the error is of type validator.ValidationErrors
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		var out strings.Builder
		out.WriteString("Validation failed with the following errors:\n")

		for _, fieldErr := range validationErrors {
			out.WriteString(fmt.Sprintf("  - Field '%s' %s \n", fieldErr.Namespace(), formatFieldError(fieldErr)))
		}
		return out.String()
	}

	return err.Error()
}

// formatFieldError creates a specific message for a single validation error.
func formatFieldError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is a required field"
	case "gt":
		return fmt.Sprintf("must be greater than %s", fe.Param())
	case "min":
		return fmt.Sprintf("must have at least %s item(s)", fe.Param())
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", fe.Param())
	case "dive":
		return "had an error in one of its items"
	default:
		return fmt.Sprintf("failed the '%s' validation", fe.Tag())
	}
}
