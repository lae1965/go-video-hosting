package validator

import (
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"slices"
)

var avatarExtentions = []string{".png", ".jpeg", ".jpg", ".svg", ".gif", ".webp", ".avif"}

var videoExtensions = []string{
	".3gp", ".avi", ".flv", ".m4v", ".mkv", ".mov", ".mp4", ".mpeg", ".mpg", ".wmv", ".webm"}

type Validator struct {
	Validate *validator.Validate
}

func New() *Validator {
	validator := &Validator{
		Validate: validator.New(),
	}
	validator.Validate.RegisterValidation("password", PasswordValidator)
	validator.Validate.RegisterValidation("avatar", fileExtentionValidator(avatarExtentions))
	validator.Validate.RegisterValidation("videofile", fileExtentionValidator(videoExtensions))
	validator.Validate.RegisterValidation("hash", VideoHashExtentionValidator)
	validator.Validate.RegisterValidation("id_list_len_1", idListValidator(1))
	validator.Validate.RegisterValidation("id_list_len_2", idListValidator(2))
	validator.Validate.RegisterValidation("id_list_len_3", idListValidator(3))
	validator.Validate.RegisterValidation("id_list_len_4", idListValidator(4))

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

func fileExtentionValidator(extentionList []string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		fileName := fl.Field().Interface().(*multipart.FileHeader).Filename
		return isValidExtention(fileName, extentionList)
	}
}

func VideoHashExtentionValidator(fl validator.FieldLevel) bool {
	return isValidExtention(fl.Field().String(), videoExtensions)

}

func isValidExtention(fileName string, extentionList []string) bool {
	return slices.Contains(extentionList, strings.ToLower(filepath.Ext(fileName)))
}
