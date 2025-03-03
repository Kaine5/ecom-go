package validator

import (
	"fmt"
	"reflect"
	"strings"

	"ecom-go/pkg/errors"
	"github.com/go-playground/validator/v10"
)

// Validator interface defines methods for validating requests
type Validator interface {
	// Validate validates the given struct and returns validation errors
	Validate(interface{}) error
}

// CustomValidator implements the Validator interface
type CustomValidator struct {
	validator *validator.Validate
}

// NewValidator creates a new validator
func NewValidator() Validator {
	v := validator.New()

	// Register custom validation tags if needed

	// Use struct field name as the error field
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return fld.Name
		}
		return name
	})

	return &CustomValidator{
		validator: v,
	}
}

// Validate validates the given struct based on the validator tags
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Convert validation errors to application errors
		var errorItems []errors.ErrorItem

		for _, err := range err.(validator.ValidationErrors) {
			errorItems = append(errorItems, errors.ErrorItem{
				Field:   err.Field(),
				Message: formatValidationError(err),
				Value:   err.Value(),
			})
		}

		return errors.WithErrors(
			errors.NewBadRequestError("validation failed", err),
			errorItems,
		)
	}

	return nil
}

// formatValidationError formats a validation error into a human-readable message
func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Must be a valid email address"
	case "min":
		return fmt.Sprintf("Must be at least %s characters long", err.Param())
	case "max":
		return fmt.Sprintf("Must not be longer than %s characters", err.Param())
	case "oneof":
		return fmt.Sprintf("Must be one of: %s", err.Param())
	default:
		return fmt.Sprintf("Failed %s validation", err.Tag())
	}
}
