package database

import (
	"database/sql"
	"fmt"
	"go-video-hosting/internal/errors"
	"go-video-hosting/internal/model"
	"strings"

	"github.com/jmoiron/sqlx"
)

type ChannelPostgres struct {
	dbSql *sqlx.DB
}

func NewChannelPostgres(dbSql *sqlx.DB) *ChannelPostgres {
	return &ChannelPostgres{dbSql: dbSql}
}

func (cp *ChannelPostgres) IsUserExist(userId int) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM USERS WHERE id = $1)"

	var exist bool
	if err := cp.dbSql.Get(&exist, query, userId); err != nil {
		return false, err
	}

	return exist, nil
}

func (cp *ChannelPostgres) IsTitlelUniqueForUser(userId int, title string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM CHANNEL WHERE title = $1 AND userId = $2)"

	var exist bool
	if err := cp.dbSql.Get(&exist, query, title, userId); err != nil {
		return false, err
	}

	return !exist, nil
}

func (cp *ChannelPostgres) CreateChannel(transaction *sqlx.Tx, userId int, title string, description string) (int, *errors.AppError) {
	query := "INSERT INTO CHANNEL (userId, title, description) VALUES ($1, $2, $3) RETURNING id"

	row := transaction.QueryRow(query, userId, title, description)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, errors.New(errors.UnknownError, err.Error())
	}

	return id, nil
}

func (cp *ChannelPostgres) UpdateChannel(userId int, channelId int, data map[string]string) *errors.AppError {
	clauses := []string{}
	args := []interface{}{}
	i := 1
	for key, value := range data {
		clauses = append(clauses, fmt.Sprintf("%s = $%d", key, i))
		args = append(args, value)
		i++
	}

	args = append(args, channelId)
	args = append(args, userId)

	query := fmt.Sprintf("UPDATE CHANNEL SET %s WHERE id = $%d AND userId = $%d", strings.Join(clauses, ", "), i, i+1)

	result, err := cp.dbSql.Exec(query, args...)
	if err != nil {
		return errors.New(errors.UnknownError, err.Error())
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return errors.New(errors.NotFound, fmt.Sprintf("channel with Id = %d not exist", channelId))
	}

	return nil
}

func (cp *ChannelPostgres) DeleteChannel(transaction *sqlx.Tx, channelId int) (int, *errors.AppError) {
	query := "DELETE FROM CHANNEL WHERE id = $1 RETURNING userId"
	var userId int

	row := transaction.QueryRow(query, channelId)
	if err := row.Scan(&userId); err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New(errors.NotFound, fmt.Sprintf("channel with Id = %d not exist", channelId))
		}
		return 0, errors.New(errors.UnknownError, err.Error())
	}

	return userId, nil
}

func (cp *ChannelPostgres) ToggleSubscribe(transaction *sqlx.Tx, userId, channelId int) (bool, *errors.AppError) {
	isUserExist, err := cp.IsUserExist(userId)
	if err != nil {
		return false, errors.New(errors.UnknownError, err.Error())
	}
	if !isUserExist {
		return false, errors.New(errors.NotFound, fmt.Sprintf("user with userId = %d not exist", userId))
	}

	query := "SELECT EXISTS (SELECT 1 FROM SUBSCRIBER WHERE userId = $1 AND channelId = $2)"

	var exist bool
	if err := cp.dbSql.Get(&exist, query, userId, channelId); err != nil {
		return false, errors.New(errors.UnknownError, err.Error())
	}

	if exist {
		query = "DELETE FROM SUBSCRIBER WHERE userId = $1 AND channelId = $2"
	} else {
		query = "INSERT INTO SUBSCRIBER (userId, channelId) VALUES ($1, $2)"
	}

	result, err := transaction.Exec(query, userId, channelId)
	if err != nil {
		return false, errors.New(errors.UnknownError, err.Error())
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return false, errors.New(errors.NotFound, fmt.Sprintf("channel with Id = %d not exist", channelId))
	}

	return !exist, nil
}

func (cp *ChannelPostgres) ChangeSubscribersCount(transaction *sqlx.Tx, channelId int, isNegative bool) (int, *errors.AppError) {
	delta := 1
	if isNegative {
		delta = -1
	}

	query := "UPDATE CHANNEL SET subscribersCount = subscribersCount + $1 WHERE id = $2 RETURNING subscribersCount"
	row := transaction.QueryRow(query, delta, channelId)

	var subscribersCount int
	if err := row.Scan(&subscribersCount); err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New(errors.NotFound, fmt.Sprintf("channel with Id = %d not exist", channelId))
		}
		return 0, errors.New(errors.UnknownError, err.Error())
	}

	return subscribersCount, nil
}

func (cp *ChannelPostgres) GetChannelById(channelId int) (*model.GetChannelFromDB, *errors.AppError) {
	query := "SELECT * FROM CHANNEL WHERE id = $1"

	var result model.GetChannelFromDB
	if err := cp.dbSql.Get(&result, query, channelId); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errors.NotFound, fmt.Sprintf("channel with Id = %d not exist", channelId))
		}
		return nil, errors.New(errors.UnknownError, err.Error())
	}

	return &result, nil
}

func (cp *ChannelPostgres) IsSubscribe(userId, channelId int) (bool, *errors.AppError) {
	query := "SELECT EXISTS (SELECT 1 FROM SUBSCRIBER WHERE userId = $1 AND channelId = $2)"

	var isExist bool
	if err := cp.dbSql.Get(&isExist, query, userId, channelId); err != nil {
		return false, errors.New(errors.UnknownError, err.Error())
	}

	return isExist, nil
}

func (cp *ChannelPostgres) GetAllChannelsOfUser(userId int) ([]*model.GetChannelFromDB, *errors.AppError) {
	query := "SELECT * FROM CHANNEL WHERE userId = $1"

	var channels []*model.GetChannelFromDB
	if err := cp.dbSql.Select(&channels, query, userId); err != nil {
		return nil, errors.New(errors.UnknownError, err.Error())
	}

	return channels, nil
}

func (cp *ChannelPostgres) GetSubscribingChannelsOfUser(userId int) ([]*model.SubscribeRequest, *errors.AppError) {
	query := `
		SELECT 
			channel.id AS channelId, 
			channel.userId AS userId 
		FROM 
			SUBSCRIBER 
		JOIN
			CHANNEL ON subscriber.channelId = channel.id
		WHERE subscriber.userId = $1
	`

	var response []*model.SubscribeRequest
	if err := cp.dbSql.Select(&response, query, userId); err != nil {
		return nil, errors.New(errors.UnknownError, err.Error())
	}
	if len(response) == 0 {
		return nil, errors.New(errors.EmptyField, fmt.Sprintf("user with Id = %d has not subscription to channels", userId))
	}

	return response, nil
}
