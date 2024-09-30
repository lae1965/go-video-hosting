package database

import (
	"go-video-hosting/pkg/model"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type Users interface {
	CreateUser(user model.Users) (int, error)
}

type Database struct {
	Users
}

func NewDatabase(db *sqlx.DB) *Database {
	return &Database{
		Users: NewUserPostgres(db),
	}
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
