package response

import (
	"errors"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type FieldError struct {
	Field   string `json:"field"`
	Rule    string `json:"rule"`
	Message string `json:"message"`
}

func ValidationErrors(err error) interface{} {
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors

	if errors.As(err, &validationErrors) {
		formattedErrors := make([]FieldError, 0, len(validationErrors))

		for _, fieldError := range validationErrors {
			fieldName := toJSONFieldName(fieldError.Field())

			formattedErrors = append(formattedErrors, FieldError{
				Field:   fieldName,
				Rule:    fieldError.Tag(),
				Message: buildValidationMessage(fieldName, fieldError),
			})
		}

		return formattedErrors
	}

	return []FieldError{
		{
			Field:   "request",
			Rule:    "invalid",
			Message: "request body or query parameters are invalid",
		},
	}
}

func toJSONFieldName(field string) string {
	if field == "" {
		return field
	}

	runes := []rune(field)
	runes[0] = unicode.ToLower(runes[0])

	return string(runes)
}

func buildValidationMessage(field string, fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return field + " is required"

	case "email":
		return field + " must be a valid email address"

	case "min":
		return field + " must be at least " + fieldError.Param() + " characters"

	case "max":
		return field + " must be at most " + fieldError.Param() + " characters"

	case "gte":
		return field + " must be greater than or equal to " + fieldError.Param()

	case "lte":
		return field + " must be less than or equal to " + fieldError.Param()

	default:
		rule := strings.ReplaceAll(fieldError.Tag(), "_", " ")
		return field + " is invalid for rule: " + rule
	}
}
