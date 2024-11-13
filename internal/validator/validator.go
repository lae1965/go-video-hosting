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
	validator.Validate.RegisterValidation("channelIdList", ChannelIdListValidator)
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

func ChannelIdListValidator(fl validator.FieldLevel) bool {
	idList := fl.Field().String()

	idListArr := strings.Split(idList, "_")
	if len(idListArr) != 2 {
		return false
	}

	for _, idStr := range idListArr {
		if _, err := strconv.ParseInt(idStr, 10, 0); err != nil {
			return false
		}
	}

	return true
}
