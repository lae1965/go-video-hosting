package database

import (
	"database/sql"
	"go-video-hosting/pkg/model"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type Users interface {
	CreateUser(transaction *sql.Tx, user model.Users) (int, error)
}

type Token interface {
	CreateToken(transaction *sql.Tx, token model.Token) (int, error)
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
