package database

import (
	"go-video-hosting/internal/errors"
	"go-video-hosting/internal/model"
	"os"

	"cnb.cool/ordermap/ordermap"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type Users interface {
	CreateUser(transaction *sqlx.Tx, user model.Users) (int, *errors.AppError)
	GetUserByEmail(email string) (*model.Users, error)
	GetUserForRefreshById(id int) (*model.Users, error)
	GetAvatarByUserId(userId int) (string, *errors.AppError)
	UpdateUser(id int, data *ordermap.OrderMap) *errors.AppError
	DeleteUser(id int) *errors.AppError
	GetUserByActivateLink(activateLink string) (int, *errors.AppError)
	GetAll() ([]*model.FindUsers, error)
	GetById(id int) (*model.FindUsers, *errors.AppError)
	GetNickNameById(id int) (string, *errors.AppError)
	CheckIsUnique(key string, value string) (bool, error)
	GetPasswordByUserId(userId int) (string, *errors.AppError)
	ChangeChannelsCountOfUser(transaction *sqlx.Tx, userId int, isIncrement bool) *errors.AppError
}

type Token interface {
	CreateToken(transaction *sqlx.Tx, token model.Token) (int, error)
	UpdateToken(tokenId int, token string) error
	RemoveToken(tokenId int) error
	GetTokenIdByToken(token string) (int, error)
	DeleteTokenFromOtherDevices(userId int, refreshTokenId int) error
}

type Channel interface {
	IsUserExist(userId int) (bool, error)
	IsTitlelUniqueForUser(userId int, title string) (bool, error)
	CreateChannel(transaction *sqlx.Tx, userId int, title string, description string) (int, *errors.AppError)
	UpdateChannel(userId int, channelId int, data map[string]string) *errors.AppError
	DeleteChannel(transaction *sqlx.Tx, channelId int) (int, *errors.AppError)
	ToggleSubscribe(transaction *sqlx.Tx, userId, channelId int) (bool, *errors.AppError)
	ChangeSubscribersCount(transaction *sqlx.Tx, channelId int, isNegative bool) (int, *errors.AppError)
	GetChannelById(channelId int) (*model.GetChannelFromDB, *errors.AppError)
	IsSubscribe(userId, channelId int) (bool, *errors.AppError)
	GetAllChannelsOfUser(userId int) ([]*model.GetChannelFromDB, *errors.AppError)
	GetSubscribingChannelsOfUser(userId int) ([]*model.SubscribeRequest, *errors.AppError)
}

type Playlist interface {
	IsChannelExist(channelId int) (bool, error)
	IsTitlelUniqueForChannel(channelId int, title string) (bool, error)
	CreatePlaylist(transaction *sqlx.Tx, channelId int, title string, description string) (int, *errors.AppError)
	UpdatePlaylist(transaction *sqlx.Tx, channelId int, playlistId int, data map[string]string) *errors.AppError
	DeletePlaylist(playlisyId int) *errors.AppError
	GetPlaylistById(playlistId int) (*model.GetPlaylist, *errors.AppError)
	GetAllPlaylistsOfChannel(channelId int) ([]*model.GetPlaylist, *errors.AppError)
}

type Video interface {
	IsTitleUniqueForChannel(channelId int, title string) (bool, error)
	CreateNewVideo(transaction *sqlx.Tx, video *model.CreateVideo) (int, *errors.AppError)
	UpdateVideoInfo(transaction *sqlx.Tx, id int, data *ordermap.OrderMap) *errors.AppError
	GetVideoHashNameById(id int) (string, *errors.AppError)
	IncrementViewsCount(videoId int) error
}

type VideoHistoty interface {
	CreateVideoHistory(userId, videoId int) error
}

type Database struct {
	Users
	Token
	Channel
	Playlist
	Video
	VideoHistoty
	dbSql *sqlx.DB
}

func New(dbSql *sqlx.DB) *Database {
	return &Database{
		Users:        NewUserPostgres(dbSql),
		Token:        NewTokenPostgres(dbSql),
		Channel:      NewChannelPostgres(dbSql),
		Playlist:     NewPlaylistPostgres(dbSql),
		Video:        NewVideoPostgres(dbSql),
		VideoHistoty: NewVideoHistoryPostgres(dbSql),
		dbSql:        dbSql,
	}
}

func (db *Database) BeginTransaction() (*sqlx.Tx, error) {
	return db.dbSql.Beginx()
}

func Connection() (*sqlx.DB, error) {
	return NewPostgresDB(Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
}
