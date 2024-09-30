package service

import (
	"crypto/sha1"
	"fmt"
	"go-video-hosting/pkg/database"
	"go-video-hosting/pkg/model"
	"os"
)

type UserService struct {
	repo database.Users
}

func NewUserService(repo database.Users) *UserService {
	return &UserService{repo: repo}
}

func (userService *UserService) CreateUser(user model.Users) (int, error) {
	user.Password = userService.GenerateHashPassword(user.Password)

	return userService.repo.CreateUser(user)
}

func (userService *UserService) GenerateHashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(os.Getenv("SALT"))))
}
