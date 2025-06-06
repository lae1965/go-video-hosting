package validator

import (
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	Validate *validator.Validate
}

func New() *Validator {
	validator := &Validator{
		Validate: validator.New(),
	}
	validator.Validate.RegisterValidation("password", PasswordValidator)
	validator.Validate.RegisterValidation("avatar", AvatarValidator)
	validator.Validate.RegisterValidation("id_list_len_1", idListValidator(1))
	validator.Validate.RegisterValidation("id_list_len_2", idListValidator(2))
	validator.Validate.RegisterValidation("id_list_len_3", idListValidator(3))

	return validator
}

var fileExtentions = [7]string{".png", ".jpeg", ".jpg", ".svg", ".gif", ".webp", ".avif"}

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

func AvatarValidator(fl validator.FieldLevel) bool {
	fileName := fl.Field().Interface().(*multipart.FileHeader).Filename
	ext := strings.ToLower(filepath.Ext(fileName))

	for _, extention := range fileExtentions {
		if extention == ext {
			return true
		}
	}

	return false
}

func idListValidator(idsCount int) validator.Func {
	return func(fl validator.FieldLevel) bool {
		idListArr := strings.Split(fl.Field().String(), "_")
		if len(idListArr) < idsCount {
			return false
		}

		for _, idStr := range idListArr {
			if _, err := strconv.ParseInt(idStr, 10, 0); err != nil {
				return false
			}
		}

		return true
	}
}
