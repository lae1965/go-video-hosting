package database

import (
	"database/sql"
	"go-video-hosting/pkg/model"

	"github.com/jmoiron/sqlx"
)

type UserPosrgres struct {
	dbSql *sqlx.DB
}

func NewUserPostgres(dbSql *sqlx.DB) *UserPosrgres {
	return &UserPosrgres{dbSql: dbSql}
}

func (userPostgres *UserPosrgres) CreateUser(transaction *sql.Tx, user model.Users) (int, error) {
	query := "INSERT INTO USERS (nickName, email, password, activateLink) values ($1, $2, $3, $4) RETURNING id"

	row := transaction.QueryRow(query, user.NickName, user.Email, user.Password, user.ActivateLink)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
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
