package service

import (
	"database/sql"
	"go-video-hosting/gRPC/client"
	"go-video-hosting/internal/errors"
	"go-video-hosting/pkg/database"
	"go-video-hosting/pkg/model"
)

type Users interface {
	CreateUser(user model.Users) (*model.UserCreateResponse, *errors.ErrorRes)
	Login(user model.Users) (*model.UserResponse, *errors.ErrorRes)
	Logout(refreshTokenId int) error
	Refresh(refreshToken string) (*model.UserResponse, *errors.ErrorRes)
	SaveAvatar(id int, fileName string) *errors.ErrorRes
	GetAvatar(id int, sendChunk func(int64, string, []byte) error) *errors.ErrorRes
	DeleteAvatar(id int) *errors.ErrorRes
	UpdateUser(id int, data map[string]interface{}) *errors.ErrorRes
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

func NewService(db *database.Database, grpcClient grpcclient.FilesGRPCClient) *Service {
	return &Service{
		Users: NewUserService(db.Users, NewTokenService(db.Token), db.BeginTransaction, grpcClient),
		Token: NewTokenService(db.Token),
	}
}
