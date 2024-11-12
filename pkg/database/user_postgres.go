package database

import (
	"database/sql"
	"fmt"
	"go-video-hosting/internal/errors"
	"go-video-hosting/pkg/model"
	"net/http"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserPostgres struct {
	dbSql *sqlx.DB
}

func NewUserPostgres(dbSql *sqlx.DB) *UserPostgres {
	return &UserPostgres{dbSql: dbSql}
}

func (userPostgres *UserPostgres) CreateUser(transaction *sql.Tx, user model.Users) (int, *errors.ErrorRes) {
	query := "INSERT INTO USERS (nickName, email, password, activateLink) values ($1, $2, $3, $4) RETURNING id"

	row := transaction.QueryRow(query, user.NickName, user.Email, user.Password, user.ActivateLink)

	var id int
	if err := row.Scan(&id); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == UniqueViolation {
				return 0, &errors.ErrorRes{Code: http.StatusConflict, Message: fmt.Sprintf("unique violation: %s", err.Error())}
			}
		}
		return 0, &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	return id, nil
}

func (userPostgres *UserPostgres) GetUserByEmail(email string) (*model.Users, error) {
	query := "SELECT * FROM USERS WHERE email=$1"

	var user model.Users
	if err := userPostgres.dbSql.Get(&user, query, email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (userPostgres *UserPostgres) GetUserById(id int) (*model.Users, error) {
	query := "SELECT id, nickName, email, role FROM USERS WHERE id=$1"

	row := userPostgres.dbSql.QueryRow(query, id)

	var user model.Users
	if err := row.Scan(&user.Id, &user.NickName, &user.Email, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func (userPostgres *UserPostgres) GetAvatarByUserId(userId int) (string, *errors.ErrorRes) {
	query := "SELECT avatar FROM USERS WHERE id=$1"

	row := userPostgres.dbSql.QueryRow(query, userId)

	var avatar string
	if err := row.Scan(&avatar); err != nil {
		if err == sql.ErrNoRows {
			return "", &errors.ErrorRes{Code: http.StatusNotFound, Message: fmt.Sprintf("user with Id = %d not exist", userId)}
		}
		return "", &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return avatar, nil
}

func (userPostgres *UserPostgres) UpdateUser(id int, data map[string]interface{}) *errors.ErrorRes {
	clauses := []string{}
	args := []interface{}{}
	i := 1
	for key, value := range data {
		clauses = append(clauses, fmt.Sprintf("%s = $%d", key, i))
		args = append(args, value)
		i++
	}
	args = append(args, id)

	query := fmt.Sprintf("UPDATE USERS SET %s WHERE id = $%d", strings.Join(clauses, ", "), i)

	result, err := userPostgres.dbSql.Exec(query, args...)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == UniqueViolation {
				return &errors.ErrorRes{Code: http.StatusConflict, Message: fmt.Sprintf("unique violation: %s", err.Error())}
			}
		}
		return &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return &errors.ErrorRes{Code: http.StatusNotFound, Message: fmt.Sprintf("user with Id = %d not exist", id)}
	}
	return nil
}

func (userPostgres *UserPostgres) DeleteUser(id int) *errors.ErrorRes {
	query := "DELETE FROM USERS WHERE id = $1"

	result, err := userPostgres.dbSql.Exec(query, id)
	if err != nil {
		return &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return &errors.ErrorRes{Code: http.StatusNotFound, Message: fmt.Sprintf("user with Id = %d not exist", id)}
	}

	return nil
}

func (userPostgres *UserPostgres) FindUserByActivateLink(activateLink string) (int, *errors.ErrorRes) {
	query := "SELECT id FROM USERS WHERE activateLink = $1"

	var id int
	if err := userPostgres.dbSql.Get(&id, query, activateLink); err != nil {
		if err == sql.ErrNoRows {
			return 0, &errors.ErrorRes{Code: http.StatusBadRequest, Message: fmt.Sprintf("user with Id = %d not exist", id)}
		}
		return 0, &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return id, nil
}

func (userPostgres *UserPostgres) FindAll() ([]*model.FindUsers, error) {
	query := "SELECT id, nickName, email, firstName, lastName, birthDate, role, isBanned, channelsCount, createTimestamp FROM USERS"

	users := []*model.FindUsers{}
	if err := userPostgres.dbSql.Select(&users, query); err != nil {
		return nil, err
	}

	return users, nil
}

func (userPostgres *UserPostgres) FindById(id int) (*model.FindUsers, *errors.ErrorRes) {
	query := "SELECT id, nickName, email, firstName, lastName, birthDate, role, isBanned, channelsCount, createTimestamp FROM USERS WHERE id = $1"

	var user model.FindUsers
	if err := userPostgres.dbSql.Get(&user, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, &errors.ErrorRes{Code: http.StatusNotFound, Message: fmt.Sprintf("user with Id = %d not exist", id)}
		}
		return nil, &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return &user, nil
}

func (userPostgres *UserPostgres) FindNickNameById(id int) (string, *errors.ErrorRes) {
	query := "SELECT nickName FROM USERS WHERE id = $1"

	var nickName string
	if err := userPostgres.dbSql.Get(&nickName, query, id); err != nil {
		if err == sql.ErrNoRows {
			return "", &errors.ErrorRes{Code: http.StatusNotFound, Message: fmt.Sprintf("user with Id = %d not exist", id)}
		}
		return "", &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return nickName, nil
}

func (userPostgres *UserPostgres) CheckIsUnique(key string, value string) (bool, error) {
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM USERS WHERE %s = $1)", key)

	var exist bool
	if err := userPostgres.dbSql.Get(&exist, query, value); err != nil {
		return false, err
	}

	return !exist, nil
}

func (userPostgres *UserPostgres) GetPasswordByUserId(userId int) (string, *errors.ErrorRes) {
	query := "SELECT password FROM USERS WHERE id = $1"

	var password string
	if err := userPostgres.dbSql.Get(&password, query, userId); err != nil {
		if err == sql.ErrNoRows {
			return "", &errors.ErrorRes{Code: http.StatusNotFound, Message: fmt.Sprintf("user with Id = %d not exist", userId)}
		}
		return "", &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return password, nil
}

func (userPostgres *UserPostgres) ChangeChannelsCountOfUser(userId int, isIncrement bool) *errors.ErrorRes {
	query := "UPDATE USERS SET channelsCount = channelsCount + $1 WHERE id = $2"
	delta := 1
	if !isIncrement {
		delta = -1
	}

	result, err := userPostgres.dbSql.Exec(query, delta, userId)
	if err != nil {
		return &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	if row, _ := result.RowsAffected(); row == 0 {
		return &errors.ErrorRes{Code: http.StatusNotFound, Message: fmt.Sprintf("user with Id = %d not exist", userId)}
	}

	return nil
}
