package service

import (
	"database/sql"
	"go-video-hosting/pkg/database"
	"go-video-hosting/pkg/model"
)

type Users interface {
	CreateUser(user model.Users) (*model.UserCreateResponse, error)
	Login(user model.Users) (*model.UserResponse, error)
	Logout(refreshTokenId int) error
	Refresh(refreshToken string) (*model.UserResponse, error)
}

type Token interface {
	CreateTokens(transaction *sql.Tx, user model.Users, refreshTokenId int) (*model.TokenResponse, error)
	RemoveToken(tokenId int) error
	ValidateToken(tokenString string, tokenKey string) (int, error)
	GetTokenIdByToken(token string) (int, error)
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
