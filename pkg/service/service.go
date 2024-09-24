package service

import "go-video-hosting/pkg/repository"

type User interface {
}

type Service struct {
	User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{}
}
