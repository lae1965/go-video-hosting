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
	DeleteUser(id int) *errors.ErrorRes
	Activate(activateLink string) *errors.ErrorRes
	FindAll() ([]*model.FindUsers, error)
	FindById(id int) (*model.FindUsers, *errors.ErrorRes)
	FindNickNameById(id int) (string, *errors.ErrorRes)
	CheckIsNickNameEmailUnique(nickName string, email string) (bool, string, error)
	ChangePassword(userId int, refreshTokenId int, oldPassword string, newPassword string) *errors.ErrorRes
}

type Token interface {
	CreateTokens(transaction *sql.Tx, user model.Users, refreshTokenId int) (*model.TokenResponse, error)
	RemoveToken(tokenId int) error
	ValidateToken(tokenString string, tokenKey string) (int, error)
	GetTokenIdByToken(token string) (int, error)
	DeleteTokenFromOtherDevices(userId int, refreshTokenId int) error
}

type Channel interface {
	CreateChannel(userId int, title string, description string) (int, *errors.ErrorRes)
	UpdateChannel(userId int, channelId int, data map[string]string) *errors.ErrorRes
	DeleteChannel(channelId int) *errors.ErrorRes
	ToggleSubscribe(userId, channelId int) (*model.SubscribeRespose, *errors.ErrorRes)
	GetChannelById(userId, channelId int) (*model.GetChannelResponse, *errors.ErrorRes)
	GetAllChannelsOfUser(userId int) ([]*model.GetAllChannelsResponse, *errors.ErrorRes)
	GetAllIdListOfUser(userId int) ([]string, *errors.ErrorRes)
}

type Service struct {
	Users
	Token
	Channel
}

func NewService(db *database.Database, grpcClient grpcclient.FilesGRPCClient) *Service {
	return &Service{
		Users:   NewUserService(db.Users, NewTokenService(db.Token), db.BeginTransaction, grpcClient),
		Token:   NewTokenService(db.Token),
		Channel: NewChannelService(db.Channel, db.Users),
	}
}
