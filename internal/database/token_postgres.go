package database

import (
	"database/sql"
	"go-video-hosting/internal/model"

	"github.com/jmoiron/sqlx"
)

type TokenPostgres struct {
	dbSql *sqlx.DB
}

func NewTokenPostgres(dbSql *sqlx.DB) *TokenPostgres {
	return &TokenPostgres{dbSql: dbSql}
}

func (tp *TokenPostgres) CreateToken(transaction *sql.Tx, token model.Token) (int, error) {
	var row *sql.Row
	query := "INSERT INTO TOKEN (token, userId) values ($1, $2) RETURNING id"

	if transaction == nil {
		row = tp.dbSql.QueryRow(query, token.Token, token.UserId)
	} else {
		row = transaction.QueryRow(query, token.Token, token.UserId)
	}

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (tp *TokenPostgres) UpdateToken(tokenId int, token string) error {
	query := "UPDATE TOKEN SET token=$1 WHERE id=$2;"

	_, err := tp.dbSql.Exec(query, token, tokenId)

	return err
}

func (tp *TokenPostgres) RemoveToken(tokenId int) error {
	query := "DELETE FROM TOKEN WHERE id=$1"
	_, err := tp.dbSql.Exec(query, tokenId)

	return err
}

func (tp *TokenPostgres) GetTokenIdByToken(token string) (int, error) {
	query := "SELECT id FROM TOKEN WHERE token=$1"

	var id int
	err := tp.dbSql.Get(&id, query, token)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (tp *TokenPostgres) DeleteTokenFromOtherDevices(userId int, refreshTokenId int) error {
	query := "DELETE FROM TOKEN WHERE userId = $1 AND id != $2"

	_, err := tp.dbSql.Exec(query, userId, refreshTokenId)

	return err
}
