package service

import (
	"database/sql"
	"go-video-hosting/gRPC/client"
	"go-video-hosting/internal/database"
	"go-video-hosting/internal/errors"
	"go-video-hosting/internal/model"

	"cnb.cool/ordermap/ordermap"
)

type Users interface {
	CreateUser(user model.Users) (*model.UserCreateResponse, *errors.AppError)
	Login(user model.Users) (*model.UserResponse, *errors.AppError)
	Logout(refreshTokenId int) error
	Refresh(refreshToken string) (*model.UserResponse, *errors.AppError)
	SaveAvatar(id int, fileName string) *errors.AppError
	GetAvatar(id int, sendChunk func(int64, string, []byte) error) *errors.AppError
	DeleteAvatar(id int) *errors.AppError
	UpdateUser(id int, data *ordermap.OrderMap) *errors.AppError
	DeleteUser(id int) *errors.AppError
	Activate(activateLink string) *errors.AppError
	GetAll() ([]*model.FindUsers, error)
	GetById(id int) (*model.FindUsers, *errors.AppError)
	GetNickNameById(id int) (string, *errors.AppError)
	CheckIsNickNameEmailUnique(nickName string, email string) (bool, string, error)
	ChangePassword(userId int, refreshTokenId int, oldPassword string, newPassword string) *errors.AppError
}

type Token interface {
	CreateTokens(transaction *sql.Tx, user model.Users, refreshTokenId int) (*model.TokenResponse, error)
	RemoveToken(tokenId int) error
	ValidateToken(tokenString string, tokenKey string) (int, error)
	GetTokenIdByToken(token string) (int, error)
	DeleteTokenFromOtherDevices(userId int, refreshTokenId int) error
}

type Channel interface {
	CreateChannel(userId int, title string, description string) (int, *errors.AppError)
	UpdateChannel(userId int, channelId int, data map[string]string) *errors.AppError
	DeleteChannel(channelId int) *errors.AppError
	ToggleSubscribe(userId, channelId int) (*model.SubscribeRespose, *errors.AppError)
	GetChannelById(userId, channelId int) (*model.GetChannelResponse, *errors.AppError)
	GetAllChannelsOfUser(userId int) ([]*model.GetAllChannelsResponse, *errors.AppError)
	GetAllIdListOfUser(userId int) ([]string, *errors.AppError)
}

type Service struct {
	Users
	Token
	Channel
}

func New(db *database.Database, grpcClient grpcclient.FilesGRPCClient) *Service {
	return &Service{
		Users:   NewUserService(db.Users, NewTokenService(db.Token), db.BeginTransaction, grpcClient),
		Token:   NewTokenService(db.Token),
		Channel: NewChannelService(db.Channel, db.Users, db.BeginTransaction),
	}
}
