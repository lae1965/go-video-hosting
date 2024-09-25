package repository

import (
	"go-video-hosting/pkg/model"

	"github.com/jmoiron/sqlx"
)

type UserPosrgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPosrgres {
	return &UserPosrgres{db: db}
}

func (userPostgres *UserPosrgres) CreateUser(user model.Users) (int, error) {
	row := userPostgres.db.QueryRow("INSERT INTO USERS (nickName, email, passwordHash) values ($1, $2, $3) RETURNING id", user.NickName, user.Email, user.Password)

	var id int
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}
