package service

import (
	"crypto/sha1"
	"fmt"
	"go-video-hosting/pkg/model"
	"go-video-hosting/pkg/repository"
)

const salt = "sdaf54jfbjbjyjb6bnaSldHNVV8d0qwjeh"

type UserService struct {
	repo repository.Users
}

func NewUserService(repo repository.Users) *UserService {
	return &UserService{repo: repo}
}

func (userService *UserService) CreateUser(user model.Users) (int, error) {
	user.Password = userService.GenerateHashPassword(user.Password)

	return userService.repo.CreateUser(user)
}

func (userService *UserService) GenerateHashPassword(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
