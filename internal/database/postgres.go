package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const UniqueViolation = "23505" // ошибка SQL - нарушение уникальности

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	dbSql, err := sqlx.Open("postgres", fmt.Sprintf(
		`host=%s port=%s user=%s dbname=%s password=%s sslmode=%s`,
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode,
	))

	if err != nil {
		return nil, err
	}

	if err = dbSql.Ping(); err != nil {
		return nil, err
	}

	return dbSql, nil
}
