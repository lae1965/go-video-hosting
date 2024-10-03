package model

import (
	_ "github.com/go-playground/validator/v10"

	"time"
)

type Role int

const (
	UserRole Role = iota
	AdminRole
	SuperAdminRole
)

type Users struct {
	Id              int       `json:"id"`
	NickName        string    `json:"nickName" validate:"required,min=3,max=30"`
	Email           string    `json:"email" validate:"required,email"`
	Password        string    `json:"password" validate:"password"`
	FirstName       string    `json:"firstName"`
	LastName        string    `json:"lastName"`
	BirthDate       time.Time `json:"birthDate"`
	Role            Role      `json:"role"`
	ActivateLink    string    `json:"activateLink"`
	IsActivate      bool      `json:"isActivate"`
	IsBanned        bool      `json:"isBanned"`
	ChannelsCount   int       `json:"channelsCount"`
	CreateTimestamp time.Time `json:"createTimestamp"`
}

type UserResponse struct {
	*TokenResponse
	UserId int
}
