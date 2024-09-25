package service

import (
	"go-video-hosting/pkg/model"
	"go-video-hosting/pkg/repository"
)

type Users interface {
	CreateUser(user model.Users) (int, error)
}

type Service struct {
	Users
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Users: NewUserService(repo.Users),
	}
}
