package database

import (
	"database/sql"
	"fmt"
	"go-video-hosting/internal/errors"
	"go-video-hosting/internal/model"
	"strings"

	"cnb.cool/ordermap/ordermap"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserPostgres struct {
	dbSql *sqlx.DB
}

func NewUserPostgres(dbSql *sqlx.DB) *UserPostgres {
	return &UserPostgres{dbSql: dbSql}
}

func (userPostgres *UserPostgres) CreateUser(transaction *sql.Tx, user model.Users) (int, *errors.AppError) {
	query := "INSERT INTO USERS (nickName, email, password, activateLink) values ($1, $2, $3, $4) RETURNING id"

	row := transaction.QueryRow(query, user.NickName, user.Email, user.Password, user.ActivateLink)

	var id int
	if err := row.Scan(&id); err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == UniqueViolation {
				return 0, errors.New(errors.NotUnique, fmt.Sprintf("unique violation: %s", err.Error()))
			}
		}
		return 0, errors.New(errors.UnknownError, err.Error())
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

func (userPostgres *UserPostgres) GetUserForRefreshById(id int) (*model.Users, error) {
	query := "SELECT id, nickName, email, role FROM USERS WHERE id=$1"

	row := userPostgres.dbSql.QueryRow(query, id)

	var user model.Users
	if err := row.Scan(&user.Id, &user.NickName, &user.Email, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func (userPostgres *UserPostgres) GetAvatarByUserId(userId int) (string, *errors.AppError) {
	query := "SELECT avatar FROM USERS WHERE id=$1"

	row := userPostgres.dbSql.QueryRow(query, userId)

	var avatar string
	if err := row.Scan(&avatar); err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New(errors.NotFound, fmt.Sprintf("user with Id = %d not exist", userId))
		}
		return "", errors.New(errors.UnknownError, err.Error())
	}

	return avatar, nil
}

func (userPostgres *UserPostgres) UpdateUser(id int, data *ordermap.OrderMap) *errors.AppError {
	clauses := []string{}
	args := []interface{}{}
	i := 1

	data.Range(func(key, value interface{}) bool {
		clauses = append(clauses, fmt.Sprintf("%s = $%d", key, i))
		args = append(args, value)
		i++
		return true
	})
	args = append(args, id)

	query := fmt.Sprintf("UPDATE USERS SET %s WHERE id = $%d", strings.Join(clauses, ", "), i)

	result, err := userPostgres.dbSql.Exec(query, args...)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == UniqueViolation {
				return errors.New(errors.NotUnique, fmt.Sprintf("unique violation: %s", err.Error()))
			}
		}
		return errors.New(errors.UnknownError, err.Error())
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return errors.New(errors.NotFound, fmt.Sprintf("user with Id = %d not exist", id))
	}
	return nil
}

func (userPostgres *UserPostgres) DeleteUser(id int) *errors.AppError {
	query := "DELETE FROM USERS WHERE id = $1"

	result, err := userPostgres.dbSql.Exec(query, id)
	if err != nil {
		return errors.New(errors.UnknownError, err.Error())
	}

	if row, _ := result.RowsAffected(); row == 0 {
		return errors.New(errors.NotFound, fmt.Sprintf("user with Id = %d not exist", id))
	}

	return nil
}

func (userPostgres *UserPostgres) GetUserByActivateLink(activateLink string) (int, *errors.AppError) {
	query := "SELECT id FROM USERS WHERE activateLink = $1"

	var id int
	if err := userPostgres.dbSql.Get(&id, query, activateLink); err != nil {
		if err == sql.ErrNoRows {
			return 0, errors.New(errors.NotFound, fmt.Sprintf("user with activateLink = %s not exist", activateLink))
		}
		return 0, errors.New(errors.UnknownError, err.Error())
	}

	return id, nil
}

func (userPostgres *UserPostgres) GetAll() ([]*model.FindUsers, error) {
	query := "SELECT id, nickName, email, firstName, lastName, birthDate, role, isBanned, channelsCount, createTimestamp FROM USERS"

	users := []*model.FindUsers{}
	if err := userPostgres.dbSql.Select(&users, query); err != nil {
		return nil, err
	}

	return users, nil
}

func (userPostgres *UserPostgres) GetById(id int) (*model.FindUsers, *errors.AppError) {
	query := "SELECT id, nickName, email, firstName, lastName, birthDate, role, isBanned, channelsCount, createTimestamp FROM USERS WHERE id = $1"

	var user model.FindUsers
	if err := userPostgres.dbSql.Get(&user, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(errors.NotFound, fmt.Sprintf("user with Id = %d not exist", id))
		}
		return nil, errors.New(errors.UnknownError, err.Error())
	}

	return &user, nil
}

func (userPostgres *UserPostgres) GetNickNameById(id int) (string, *errors.AppError) {
	query := "SELECT nickName FROM USERS WHERE id = $1"

	var nickName string
	if err := userPostgres.dbSql.Get(&nickName, query, id); err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New(errors.NotFound, fmt.Sprintf("user with Id = %d not exist", id))
		}
		return "", errors.New(errors.UnknownError, err.Error())
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

func (userPostgres *UserPostgres) GetPasswordByUserId(userId int) (string, *errors.AppError) {
	query := "SELECT password FROM USERS WHERE id = $1"

	var password string
	if err := userPostgres.dbSql.Get(&password, query, userId); err != nil {
		if err == sql.ErrNoRows {
			return "", errors.New(errors.NotFound, fmt.Sprintf("user with Id = %d not exist", userId))
		}
		return "", errors.New(errors.UnknownError, err.Error())
	}

	return password, nil
}

func (userPostgres *UserPostgres) ChangeChannelsCountOfUser(transaction *sql.Tx, userId int, isIncrement bool) *errors.AppError {
	query := "UPDATE USERS SET channelsCount = channelsCount + $1 WHERE id = $2"
	delta := 1
	if !isIncrement {
		delta = -1
	}

	result, err := transaction.Exec(query, delta, userId)
	if err != nil {
		return errors.New(errors.UnknownError, err.Error())
	}
	if row, _ := result.RowsAffected(); row == 0 {
		return errors.New(errors.NotFound, fmt.Sprintf("user with Id = %d not exist", userId))
	}

	return nil
}
