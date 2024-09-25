package repository

import (
	"go-video-hosting/pkg/model"

	"github.com/jmoiron/sqlx"
)

type Users interface {
	CreateUser(user model.Users) (int, error)
}

type Repository struct {
	Users
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Users: NewUserPostgres(db),
	}
}
