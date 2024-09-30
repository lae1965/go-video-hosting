package service

import (
	"go-video-hosting/pkg/database"
	"go-video-hosting/pkg/model"
)

type Users interface {
	CreateUser(user model.Users) (int, error)
}

type Service struct {
	Users
}

func NewService(repo *database.Database) *Service {
	return &Service{
		Users: NewUserService(repo.Users),
	}
}
