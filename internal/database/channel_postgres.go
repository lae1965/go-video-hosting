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

func (channelPostgres *ChannelPostgres) IsUserExist(userId int) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM USERS WHERE id = $1)"

	var exist bool
	if err := channelPostgres.dbSql.Get(&exist, query, userId); err != nil {
		return false, err
	}

	return exist, nil
}

func (channelPostgres *ChannelPostgres) IsTitlelUniqueForUser(userId int, title string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM CHANNEL WHERE title = $1 AND userId = $2)"

	var exist bool
	if err := channelPostgres.dbSql.Get(&exist, query, title, userId); err != nil {
		return false, err
	}

	return !exist, nil
}

func (channelPostgres *ChannelPostgres) CreateChannel(transaction *sql.Tx, userId int, title string, description string) (int, *errors.AppError) {
	isUserExist, err := channelPostgres.IsUserExist(userId)
	if err != nil {
		return 0, errors.New(errors.UnknownError, err.Error())
	}
	if !isUserExist {
		return 0, errors.New(errors.NotFound, fmt.Sprintf("user with userId = %d not exist", userId))
	}

	isUnique, err := channelPostgres.IsTitlelUniqueForUser(userId, title)
	if err != nil {
		return 0, errors.New(errors.UnknownError, err.Error())
	}
	if !isUnique {
		return 0, errors.New(errors.NotUnique, "user's channel name must be unique")
	}

	query := "INSERT INTO CHANNEL (userId, title, description) VALUES ($1, $2, $3) RETURNING id"

	row := transaction.QueryRow(query, userId, title, description)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, errors.New(errors.UnknownError, err.Error())
	}

	return id, nil
}

func (channelPostgres *ChannelPostgres) UpdateChannel(userId int, channelId int, data map[string]string) *errors.AppError {
	isUserExist, err := channelPostgres.IsUserExist(userId)
	if err != nil {
		return errors.New(errors.UnknownError, err.Error())
	}
	if !isUserExist {
		return errors.New(errors.NotFound, fmt.Sprintf("user with userId = %d not exist", userId))
	}

	title := ""
	clauses := []string{}
	args := []interface{}{}
	i := 1
	for key, value := range data {
		if key == "title" {
			title = value
		}
		clauses = append(clauses, fmt.Sprintf("%s = $%d", key, i))
		args = append(args, value)
		i++
	}

	if title != "" {
		isUnique, err := channelPostgres.IsTitlelUniqueForUser(userId, title)
		if err != nil {
			return errors.New(errors.UnknownError, err.Error())
		}
		if !isUnique {
			return errors.New(errors.NotUnique, "user's channel name must be unique")
		}
	}

	args = append(args, channelId)
	args = append(args, userId)

	query := fmt.Sprintf("UPDATE CHANNEL SET %s WHERE id = $%d AND userId = $%d", strings.Join(clauses, ", "), i, i+1)

	result, err := channelPostgres.dbSql.Exec(query, args...)
	if err != nil {
		return errors.New(errors.UnknownError, err.Error())
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return errors.New(errors.NotFound, fmt.Sprintf("channel with Id = %d not exist", channelId))
	}

	return nil
}

func (channelPostgres *ChannelPostgres) DeleteChannel(transaction *sql.Tx, channelId int) (int, *errors.AppError) {
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

func (channelPostgres *ChannelPostgres) ToggleSubscribe(transaction *sql.Tx, userId, channelId int) (bool, *errors.AppError) {
	isUserExist, err := channelPostgres.IsUserExist(userId)
	if err != nil {
		return false, errors.New(errors.UnknownError, err.Error())
	}
	if !isUserExist {
		return false, errors.New(errors.NotFound, fmt.Sprintf("user with userId = %d not exist", userId))
	}

	query := "SELECT EXISTS (SELECT 1 FROM SUBSCRIBER WHERE userId = $1 AND channelId = $2)"

	var exist bool
	if err := channelPostgres.dbSql.Get(&exist, query, userId, channelId); err != nil {
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

func (channelPostgres *ChannelPostgres) ChangeSubscribersCount(transaction *sql.Tx, channelId int, isNegative bool) (int, *errors.AppError) {
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

func (channelPostgres *ChannelPostgres) GetChannelById(channelId int) (*model.GetChannelFromDB, *errors.AppError) {
	query := "SELECT * FROM CHANNEL WHERE id = $1"

	var result model.GetChannelFromDB
	if err := channelPostgres.dbSql.Get(&result, query, channelId); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errors.NotFound, fmt.Sprintf("channel with Id = %d not exist", channelId))
		}
		return nil, errors.New(errors.UnknownError, err.Error())
	}

	return &result, nil
}

func (channelPostgres *ChannelPostgres) IsSubscribe(userId, channelId int) (bool, *errors.AppError) {
	query := "SELECT EXISTS (SELECT 1 FROM SUBSCRIBER WHERE userId = $1 AND channelId = $2)"

	var isExist bool
	if err := channelPostgres.dbSql.Get(&isExist, query, userId, channelId); err != nil {
		return false, errors.New(errors.UnknownError, err.Error())
	}

	return isExist, nil
}

func (channelPostgres *ChannelPostgres) GetAllChannelsOfUser(userId int) ([]*model.GetChannelFromDB, *errors.AppError) {
	query := "SELECT * FROM CHANNEL WHERE userId = $1"

	var channels []*model.GetChannelFromDB
	if err := channelPostgres.dbSql.Select(&channels, query, userId); err != nil {
		return nil, errors.New(errors.UnknownError, err.Error())
	}
	if len(channels) == 0 {
		return nil, errors.New(errors.EmptyField, fmt.Sprintf("user with Id = %d has no channels", userId))
	}

	return channels, nil
}

func (channelPostgres *ChannelPostgres) GetSubscribingChannelsOfUser(userId int) ([]*model.SubscribeRequest, *errors.AppError) {
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
	if err := channelPostgres.dbSql.Select(&response, query, userId); err != nil {
		return nil, errors.New(errors.UnknownError, err.Error())
	}
	if len(response) == 0 {
		return nil, errors.New(errors.EmptyField, fmt.Sprintf("user with Id = %d has not subscription to channels", userId))
	}

	return response, nil
}
