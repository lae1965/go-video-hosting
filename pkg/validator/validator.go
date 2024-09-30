package validator

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validate *validator.Validate
}

func NewValidator() *Validator {
	validator := new(Validator)
	validator.Validate.RegisterValidation("password", PasswordValidator)
	return validator
}

func PasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if length := len(password); length < 8 || length > 40 {
		return false
	}

	wasUpper := false
	wasLower := false
	digitsCount := 0

	for _, char := range password {
		if unicode.IsSpace(char) {
			return false
		}
		if unicode.IsUpper(char) {
			wasUpper = true
		} else if unicode.IsLower(char) {
			wasLower = true
		} else if unicode.IsDigit(char) {
			digitsCount++
		}
	}

	return wasUpper && wasLower && digitsCount >= 2
}
