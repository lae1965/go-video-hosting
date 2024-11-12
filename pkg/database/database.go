package database

import (
	"database/sql"
	"go-video-hosting/internal/errors"
	"go-video-hosting/pkg/model"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type Users interface {
	CreateUser(transaction *sql.Tx, user model.Users) (int, *errors.ErrorRes)
	GetUserByEmail(email string) (*model.Users, error)
	GetUserById(id int) (*model.Users, error)
	GetAvatarByUserId(userId int) (string, *errors.ErrorRes)
	UpdateUser(id int, data map[string]interface{}) *errors.ErrorRes
	DeleteUser(id int) *errors.ErrorRes
	FindUserByActivateLink(activateLink string) (int, *errors.ErrorRes)
	FindAll() ([]*model.FindUsers, error)
	FindById(id int) (*model.FindUsers, *errors.ErrorRes)
	FindNickNameById(id int) (string, *errors.ErrorRes)
	CheckIsUnique(key string, value string) (bool, error)
	GetPasswordByUserId(userId int) (string, *errors.ErrorRes)
}

type Token interface {
	CreateToken(transaction *sql.Tx, token model.Token) (int, error)
	UpdateToken(tokenId int, token string) error
	RemoveToken(tokenId int) error
	GetTokenIdByToken(token string) (int, error)
	DeleteTokenFromOtherDevices(userId int, refreshTokenId int) error
}

type Database struct {
	Users
	Token
	dbSql *sqlx.DB
}

func NewDatabase(dbSql *sqlx.DB) *Database {
	return &Database{
		Users: NewUserPostgres(dbSql),
		Token: NewTokenPostgres(dbSql),
		dbSql: dbSql,
	}
}

func (db *Database) BeginTransaction() (*sql.Tx, error) {
	return db.dbSql.Begin()
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
