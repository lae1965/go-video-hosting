package database

import (
	"database/sql"
	"fmt"
	"go-video-hosting/internal/errors"
	"go-video-hosting/internal/model"
	"net/http"
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

func (channelPostgres *ChannelPostgres) CreateChannel(userId int, title string, description string) (int, *errors.ErrorRes) {
	isUserExist, err := channelPostgres.IsUserExist(userId)
	if err != nil {
		return 0, &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	if !isUserExist {
		return 0, &errors.ErrorRes{Code: http.StatusBadRequest, Message: fmt.Sprintf("user with userId = %d not exist", userId)}
	}

	isUnique, err := channelPostgres.IsTitlelUniqueForUser(userId, title)
	if err != nil {
		return 0, &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	if !isUnique {
		return 0, &errors.ErrorRes{Code: http.StatusConflict, Message: "user's channel name must be unique"}
	}

	query := "INSERT INTO CHANNEL (userId, title, description) VALUES ($1, $2, $3) RETURNING id"

	row := channelPostgres.dbSql.QueryRow(query, userId, title, description)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return id, nil
}

func (channelPostgres *ChannelPostgres) UpdateChannel(userId int, channelId int, data map[string]string) *errors.ErrorRes {
	isUserExist, err := channelPostgres.IsUserExist(userId)
	if err != nil {
		return &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	if !isUserExist {
		return &errors.ErrorRes{Code: http.StatusBadRequest, Message: fmt.Sprintf("user with userId = %d not exist", userId)}
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
			return &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
		}
		if !isUnique {
			return &errors.ErrorRes{Code: http.StatusConflict, Message: "user's channel name must be unique"}
		}
	}

	args = append(args, channelId)
	args = append(args, userId)

	query := fmt.Sprintf("UPDATE CHANNEL SET %s WHERE id = $%d AND userId = $%d", strings.Join(clauses, ", "), i, i+1)

	result, err := channelPostgres.dbSql.Exec(query, args...)
	if err != nil {
		return &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return &errors.ErrorRes{Code: http.StatusNotFound, Message: fmt.Sprintf("channel with Id = %d not exist", channelId)}
	}

	return nil
}

func (channelPostgres *ChannelPostgres) DeleteChannel(channelId int) (int, *errors.ErrorRes) {
	query := "DELETE FROM CHANNEL WHERE id = $1 RETURNING userId"
	var userId int

	if err := channelPostgres.dbSql.Get(&userId, query, channelId); err != nil {
		if err == sql.ErrNoRows {
			return 0, &errors.ErrorRes{Code: http.StatusNotFound, Message: fmt.Sprintf("channel with Id = %d not exist", channelId)}
		}
		return 0, &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return userId, nil
}

func (channelPostgres *ChannelPostgres) ToggleSubscribe(userId, channelId int) (bool, *errors.ErrorRes) {
	isUserExist, err := channelPostgres.IsUserExist(userId)
	if err != nil {
		return false, &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	if !isUserExist {
		return false, &errors.ErrorRes{Code: http.StatusNotFound, Message: fmt.Sprintf("user with userId = %d not exist", userId)}
	}

	query := "SELECT EXISTS (SELECT 1 FROM SUBSCRIBER WHERE userId = $1 AND channelId = $2)"

	var exist bool
	if err := channelPostgres.dbSql.Get(&exist, query, userId, channelId); err != nil {
		return false, &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	if exist {
		query = "DELETE FROM SUBSCRIBER WHERE userId = $1 AND channelId = $2"
	} else {
		query = "INSERT INTO SUBSCRIBER (userId, channelId) VALUES ($1, $2)"
	}

	result, err := channelPostgres.dbSql.Exec(query, userId, channelId)
	if err != nil {
		return false, &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return false, &errors.ErrorRes{Code: http.StatusNotFound, Message: fmt.Sprintf("channel with Id = %d not exist", channelId)}
	}

	return !exist, nil
}

func (channelPostgres *ChannelPostgres) ChangeSubscribersCount(channelId int, isNegative bool) (int, *errors.ErrorRes) {
	delta := 1
	if isNegative {
		delta = -1
	}

	query := "UPDATE CHANNEL SET subscribersCount = subscribersCount + $1 WHERE id = $2 RETURNING subscribersCount"
	row := channelPostgres.dbSql.QueryRow(query, delta, channelId)

	var subscribersCount int
	if err := row.Scan(&subscribersCount); err != nil {
		if err == sql.ErrNoRows {
			return 0, &errors.ErrorRes{Code: http.StatusNotFound, Message: fmt.Sprintf("channel with Id = %d not exist", channelId)}
		}
		return 0, &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return subscribersCount, nil
}

func (channelPostgres *ChannelPostgres) GetChannelById(channelId int) (*model.GetChannelFromDB, *errors.ErrorRes) {
	query := "SELECT * FROM CHANNEL WHERE id = $1"

	var result model.GetChannelFromDB
	if err := channelPostgres.dbSql.Get(&result, query, channelId); err != nil {
		if err == sql.ErrNoRows {
			return nil, &errors.ErrorRes{Code: http.StatusNotFound, Message: fmt.Sprintf("channel with Id = %d not exist", channelId)}
		}
		return nil, &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return &result, nil
}

func (channelPostgres *ChannelPostgres) IsSubscribe(userId, channelId int) (bool, *errors.ErrorRes) {
	query := "SELECT EXISTS (SELECT 1 FROM SUBSCRIBER WHERE userId = $1 AND channelId = $2)"

	var isExist bool
	if err := channelPostgres.dbSql.Get(&isExist, query, userId, channelId); err != nil {
		return false, &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return isExist, nil
}

func (channelPostgres *ChannelPostgres) GetAllChannelsOfUser(userId int) ([]*model.GetChannelFromDB, *errors.ErrorRes) {
	query := "SELECT * FROM CHANNEL WHERE userId = $1"

	var channels []*model.GetChannelFromDB
	if err := channelPostgres.dbSql.Select(&channels, query, userId); err != nil {
		return nil, &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	if len(channels) == 0 {
		return nil, &errors.ErrorRes{Code: http.StatusNoContent, Message: fmt.Sprintf("user with Id = %d has no channels", userId)}
	}

	return channels, nil
}

func (channelPostgres *ChannelPostgres) GetSubscribingChannelsOfUser(userId int) ([]*model.SubscribeRequest, *errors.ErrorRes) {
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
		return nil, &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	if len(response) == 0 {
		return nil, &errors.ErrorRes{Code: http.StatusNoContent, Message: fmt.Sprintf("user with Id = %d has not subscription to channels", userId)}
	}

	return response, nil
}
