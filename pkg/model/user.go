package model

import (
	"database/sql"

	_ "github.com/go-playground/validator/v10"

	"time"
)

type Role string

const (
	UserRole       Role = "userRole"
	AdminRole      Role = "adminRole"
	SuperAdminRole Role = "superAdminRole"
)

type Users struct {
	Id              int          `json:"id"`
	NickName        string       `json:"nickName"`
	Email           string       `json:"email" validate:"required,email"`
	Password        string       `json:"password" validate:"password"`
	FirstName       string       `json:"firstName"`
	LastName        string       `json:"lastName"`
	BirthDate       sql.NullTime `json:"birthDate"`
	Role            Role         `json:"role"`
	ActivateLink    string       `json:"activateLink"`
	IsActivate      bool         `json:"isActivate"`
	IsBanned        bool         `json:"isBanned"`
	ChannelsCount   int          `json:"channelsCount"`
	CreateTimestamp time.Time    `json:"createTimestamp"`
}

type UserCreateResponse struct {
	*TokenResponse
	UserId int
}

type UserResponse struct {
	*TokenResponse
	*Users
}
