package helper

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// FormatValidationErrors formats go-playground validation errors into our map format
func FormatValidationErrors(errs validator.ValidationErrors) map[string][]string {
	validationErrors := make(map[string][]string)
	for _, err := range errs {
		fieldName := strings.ToLower(err.Field())
		var errCode string
		switch err.Tag() {
		case "required":
			errCode = "IS_REQUIRED"
		case "email":
			errCode = "IS_INVALID"
		case "min":
			errCode = "TOO_SHORT"
		case "max":
			errCode = "TOO_LONG"
		case "lowercase":
			errCode = "MUST_LOWER"
		case "uppercase":
			errCode = "MUST_UPPER"
		case "symbol":
			errCode = "MUST_SYMBOL"
		case "number":
			errCode = "MUST_NUMBER"
		default:
			errCode = "IS_INVALID"
		}
		validationErrors[fieldName] = append(validationErrors[fieldName], errCode)
	}
	return validationErrors
}

// FormatValidationError extracts and formats validator.ValidationErrors from an error if present,
// otherwise it returns the raw error string.
func FormatValidationError(err error) interface{} {
	var validationErrs validator.ValidationErrors
	if errors.As(err, &validationErrs) {
		return FormatValidationErrors(validationErrs)
	}
	return err.Error()
}

// ValidateStruct validates a struct using go-playground/validator tags
// Returns a map of field errors if validation fails
func ValidateStruct(s interface{}) map[string][]string {
	errs := validate.Struct(s)
	if errs != nil {
		if validationErrors, ok := errs.(validator.ValidationErrors); ok {
			return FormatValidationErrors(validationErrors)
		}
	}
	return nil
}