package database

import (
	"database/sql"
	"go-video-hosting/pkg/model"

	"github.com/jmoiron/sqlx"
)

type TokenPosrgres struct {
	dbSql *sqlx.DB
}

func NewTokenPostgres(dbSql *sqlx.DB) *TokenPosrgres {
	return &TokenPosrgres{dbSql: dbSql}
}

func (tokenPostgres *TokenPosrgres) CreateToken(transaction *sql.Tx, token model.Token) (int, error) {
	var row *sql.Row
	query := "INSERT INTO TOKEN (token, userId) values ($1, $2) RETURNING id"

	if transaction == nil {
		row = tokenPostgres.dbSql.QueryRow(query, token.Token, token.UserId)
	} else {
		row = transaction.QueryRow(query, token.Token, token.UserId)
	}

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (tokenPostgres *TokenPosrgres) UpdateToken(transaction *sql.Tx, token model.Token) (int, error) {
	var row *sql.Row
	query := "INSERT INTO TOKEN (token, userId) values ($1, $2) RETURNING id"

	if transaction == nil {
		row = tokenPostgres.dbSql.QueryRow(query, token.Token, token.UserId)
	} else {
		row = transaction.QueryRow(query, token.Token, token.UserId)
	}

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
