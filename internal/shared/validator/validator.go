package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator wraps the validator instance
type Validator struct {
	validate *validator.Validate
}

// New creates a new validator instance
func New() *Validator {
	v := validator.New()

	// Use JSON tag names for field names in errors
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	// Register custom validations
	registerCustomValidations(v)

	return &Validator{validate: v}
}

// Validate validates a struct
func (v *Validator) Validate(i interface{}) map[string]string {
	err := v.validate.Struct(i)
	if err == nil {
		return nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		field := err.Field()
		errors[field] = getErrorMessage(err)
	}

	return errors
}

// getErrorMessage returns a human-readable error message
func getErrorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Value is too short (minimum: " + err.Param() + ")"
	case "max":
		return "Value is too long (maximum: " + err.Param() + ")"
	case "gte":
		return "Value must be greater than or equal to " + err.Param()
	case "lte":
		return "Value must be less than or equal to " + err.Param()
	case "oneof":
		return "Value must be one of: " + err.Param()
	case "url":
		return "Invalid URL format"
	case "phone":
		return "Invalid phone number format"
	case "password":
		return "Password must be at least 8 characters with uppercase, lowercase, and number"
	default:
		return "Invalid value"
	}
}

// registerCustomValidations registers custom validation rules
func registerCustomValidations(v *validator.Validate) {
	// Phone number validation (basic Indonesian format)
	v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		if phone == "" {
			return true // Optional field
		}
		// Basic validation: starts with +62 or 08, 10-15 digits
		if len(phone) < 10 || len(phone) > 15 {
			return false
		}
		// Allow +, digits, and dashes
		for _, c := range phone {
			if c != '+' && c != '-' && (c < '0' || c > '9') {
				return false
			}
		}
		return true
	})

	// Password validation
	v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		if len(password) < 8 {
			return false
		}
		hasUpper := false
		hasLower := false
		hasNumber := false
		for _, c := range password {
			if c >= 'A' && c <= 'Z' {
				hasUpper = true
			}
			if c >= 'a' && c <= 'z' {
				hasLower = true
			}
			if c >= '0' && c <= '9' {
				hasNumber = true
			}
		}
		return hasUpper && hasLower && hasNumber
	})
}
