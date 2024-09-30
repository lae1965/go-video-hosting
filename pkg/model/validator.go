package model

import (
	"unicode"

	"github.com/go-playground/validator/v10"
)

var ValidatorInstance *validator.Validate

func GetValidator() *validator.Validate {
	if ValidatorInstance == nil {
		ValidatorInstance = validator.New()
		ValidatorInstance.RegisterValidation("password", PasswordValidator)
	}
	return ValidatorInstance
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
