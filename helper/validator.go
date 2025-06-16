package helper

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct validates a struct using go-playground/validator tags
// Returns a map of field errors if validation fails
func ValidateStruct(s interface{}) map[string]string {
	errs := validate.Struct(s)
	if errs != nil {
		validationErrors := make(map[string]string)
		for _, err := range errs.(validator.ValidationErrors) {
			fieldName := strings.ToLower(err.Field())
			switch err.Tag() {
			case "required":
				validationErrors[fieldName] = fmt.Sprintf("%s is required", fieldName)
			case "email":
				validationErrors[fieldName] = fmt.Sprintf("%s must be a valid email address", fieldName)
			case "min":
				validationErrors[fieldName] = fmt.Sprintf("%s must be at least %s characters long", fieldName, err.Param())
			case "max":
				validationErrors[fieldName] = fmt.Sprintf("%s must be at most %s characters long", fieldName, err.Param())
			default:
				validationErrors[fieldName] = fmt.Sprintf("%s is invalid: %s", fieldName, err.Tag())
			}
		}
		return validationErrors
	}
	return nil
}