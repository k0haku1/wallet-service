package validation

import (
	"github.com/go-playground/validator/v10"
	"unicode"
)

func FormatValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, e := range validationErrs {
			field := toCamelCase(e.Field()) // WalletID → walletId
			switch e.Tag() {
			case "required":
				errors[field] = field + " is required"
			case "oneof":
				errors[field] = field + " must be one of [" + e.Param() + "]"
			case "gt":
				errors[field] = field + " must be greater than " + e.Param()
			default:
				errors[field] = "invalid value for " + field
			}
		}
	} else {
		errors["error"] = err.Error()
	}

	return errors
}

func toCamelCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}
