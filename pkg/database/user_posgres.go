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

type UserPosrgres struct {
	dbSql *sqlx.DB
}

func NewUserPostgres(dbSql *sqlx.DB) *UserPosrgres {
	return &UserPosrgres{dbSql: dbSql}
}

func (userPostgres *UserPosrgres) CreateUser(transaction *sql.Tx, user model.Users) (int, *errors.ErrorRes) {
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

func (userPosrgres *UserPosrgres) GetUserByEmail(email string) (*model.Users, error) {
	query := "SELECT * FROM USERS WHERE email=$1"

	var user model.Users
	if err := userPosrgres.dbSql.Get(&user, query, email); err != nil {
		return nil, err
	}

	return &user, nil
}

func (userPosrgres *UserPosrgres) GetUserById(id int) (*model.Users, error) {
	query := "SELECT id, nickName, email, role FROM USERS WHERE id=$1"

	row := userPosrgres.dbSql.QueryRow(query, id)

	var user model.Users
	if err := row.Scan(&user.Id, &user.NickName, &user.Email, &user.Role); err != nil {
		return nil, err
	}

	return &user, nil
}

func (userPosrgres *UserPosrgres) GetAvatarByUserId(userId int) (string, *errors.ErrorRes) {
	query := "SELECT avatar FROM USERS WHERE id=$1"

	row := userPosrgres.dbSql.QueryRow(query, userId)

	var avatar string
	if err := row.Scan(&avatar); err != nil {
		if err == sql.ErrNoRows {
			return "", &errors.ErrorRes{Code: http.StatusNotFound, Message: fmt.Sprintf("user with Id = %d not exist", userId)}
		}
		return "", &errors.ErrorRes{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return avatar, nil
}

func (userPosgres *UserPosrgres) UpdateUser(id int, data map[string]interface{}) *errors.ErrorRes {
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

	result, err := userPosgres.dbSql.Exec(query, args...)
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
