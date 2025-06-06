package database

import (
	"database/sql"
	"fmt"
	"go-video-hosting/internal/errors"
	"strings"

	"github.com/jmoiron/sqlx"
)

type PlaylistPostgres struct {
	dbSql *sqlx.DB
}

func NewPlaylistPostgres(dbSql *sqlx.DB) *PlaylistPostgres {
	return &PlaylistPostgres{dbSql: dbSql}
}

func (pp *PlaylistPostgres) IsChannelExist(channelId int) (bool, error) {
	query := "SELLECT EXISTS(SELECT 1 FROM CHANNEL WHERE id = $1)"

	var exist bool
	if err := pp.dbSql.Get(&exist, query, channelId); err != nil {
		return false, err
	}

	return exist, nil
}

func (pp *PlaylistPostgres) IsTitlelUniqueForChannel(channelId int, title string) (bool, error) {
	query := "SELECT EXIST (SELECT 1 FROM PLAYLIST WHERE title = $1 AND channelId = $2)"

	var exist bool
	if err := pp.dbSql.Get(&exist, query, title, channelId); err != nil {
		return false, err
	}

	return !exist, nil
}

func (pp *PlaylistPostgres) CreatePlaylist(transaction *sql.Tx, channelId int, title string, description string) (int, *errors.AppError) {
	query := "INSERT INTO PLAYLIST (channelId, title, description) VALUES ($1, $2, $3) RETURNING id"

	row := transaction.QueryRow(query, channelId, title, description)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, errors.New(errors.UnknownError, err.Error())
	}

	return id, nil
}

func (pp *PlaylistPostgres) UpdatePlaylist(transaction *sql.Tx, channelId int, playlistId int, data map[string]string) *errors.AppError {
	clauses := []string{}
	args := []interface{}{}
	i := 1
	for key, value := range data {
		clauses = append(clauses, fmt.Sprintf("%s = $%d", key, i))
		args = append(args, value)
		i++
	}

	args = append(args, playlistId)
	args = append(args, channelId)

	query := fmt.Sprintf("UPDATE PLAYLIST SET %s WHERE id = $%d AND channelId = $%d", strings.Join(clauses, ", "), i, i+1)

	var result sql.Result
	var err error
	if transaction == nil {
		result, err = pp.dbSql.Exec(query, args...)
	} else {
		result, err = transaction.Exec(query, args...)
	}
	if err != nil {
		return errors.New(errors.UnknownError, err.Error())
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return errors.New(errors.NotFound, fmt.Sprintf("playlist with Id = %d not exist", playlistId))
	}

	return nil
}

func (pp *PlaylistPostgres) DeletePlaylist(playlisyId int) *errors.AppError {
	query := "DELETE FROM PLAYLIST WHERE ID = $1"

	result, err := pp.dbSql.Exec(query, playlisyId)
	if err != nil {
		return errors.New(errors.UnknownError, err.Error())
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return errors.New(errors.NotFound, fmt.Sprintf("playlist with Id = %d not exist", playlisyId))
	}

	return nil
}
