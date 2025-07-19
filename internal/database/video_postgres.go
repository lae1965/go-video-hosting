package database

import (
	"database/sql"
	"fmt"
	"go-video-hosting/internal/errors"
	"go-video-hosting/internal/model"
	"strings"

	"cnb.cool/ordermap/ordermap"
	"github.com/jmoiron/sqlx"
)

type VideoPostgres struct {
	dbSql *sqlx.DB
}

func NewVideoPostgres(dbSql *sqlx.DB) *VideoPostgres {
	return &VideoPostgres{dbSql: dbSql}
}

func (vp *VideoPostgres) IsTitleUniqueForChannel(channelId int, title string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM VIDEO WHERE title = $1 AND channelId = $2)"

	var exist bool
	if err := vp.dbSql.Get(&exist, query, title, channelId); err != nil {
		return false, err
	}

	return !exist, nil
}

func (vp *VideoPostgres) CreateNewVideo(transaction *sqlx.Tx, video *model.CreateVideo) (int, *errors.AppError) {
	query := "INSERT INTO VIDEO (title, hashName, category, description, idList, thumbnail, duration) VALUES (:title, :hashName, :category, :description, :idList, :thumbnail, :duration) RETURNING id"

	row, err := transaction.NamedQuery(query, video)
	if err != nil {
		return 0, errors.New(errors.UnknownError, err.Error())
	}

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, errors.New(errors.UnknownError, err.Error())
	}

	return id, nil
}

func (vp *VideoPostgres) UpdateVideoInfo(transaction *sqlx.Tx, id int, data *ordermap.OrderMap) *errors.AppError {
	clauses := []string{}
	args := []interface{}{}
	i := 1

	data.Range(func(key, value interface{}) bool {
		clauses = append(clauses, fmt.Sprintf("%s = %d", key, i))
		args = append(args, value)
		i++
		return true
	})
	args = append(args, id)

	query := fmt.Sprintf("UPDATE VIDEO SET %s WHERE id = %d", strings.Join(clauses, ", "), i)

	var exec func(query string, args ...any) (sql.Result, error)
	if transaction == nil {
		exec = vp.dbSql.Exec
	} else {
		exec = transaction.Exec
	}

	result, err := exec(query, args...)
	if err != nil {
		return errors.New(errors.UnknownError, err.Error())
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return errors.New(errors.NotFound, fmt.Sprintf("video with Id = %d not exist", id))
	}

	return nil
}

func (vp *VideoPostgres) GetVideoHashNameById(id int) (string, *errors.AppError) {
	query := "SELECT hashName FROM VIDEO WHERE id = $1"

	var hashName string
	if err := vp.dbSql.Get(&hashName, query, id); err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New(errors.NotFound, fmt.Sprintf("video with Id = %d not exist", id))
		}
		return "", errors.New(errors.UnknownError, err.Error())
	}

	return hashName, nil
}

func (vp *VideoPostgres) IncrementViewsCount(videoId int) error {
	query := "UPDATE VIDEO SET viewsCount = viewsCount + 1 WHERE id = $1"

	_, err := vp.dbSql.Exec(query, videoId)
	return err
}
