package service

import (
	"database/sql"
	"go-video-hosting/pkg/database"
	"go-video-hosting/pkg/model"
)

type Users interface {
	CreateUser(user model.Users) (*model.UserCreateResponse, error)
	Login(user model.Users) (*model.UserResponse, error)
}

type Token interface {
	CreateTokens(transaction *sql.Tx, user model.Users) (*model.TokenResponse, error)
}

type Service struct {
	Users
	Token
}

func NewService(db *database.Database) *Service {
	return &Service{
		Users: NewUserService(db.Users, NewTokenService(db.Token), db.BeginTransaction),
		Token: NewTokenService(db.Token),
	}
}
