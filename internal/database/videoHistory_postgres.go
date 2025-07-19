package database

import "github.com/jmoiron/sqlx"

type VideoHistoryPostgres struct {
	dbSql *sqlx.DB
}

func NewVideoHistoryPostgres(dbSql *sqlx.DB) *VideoHistoryPostgres {
	return &VideoHistoryPostgres{dbSql: dbSql}
}

func (vhp *VideoHistoryPostgres) CreateVideoHistory(userId, videoId int) error {
	query := "INSERTY INTO VIDEOHISTORY (userId, videoId) VALUES ($1, $2)"

	_, err := vhp.dbSql.Query(query, userId, videoId)

	return err
}
