package model

import (
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
	Id              int       `db:"id" json:"id"`
	NickName        string    `db:"nickname" json:"nickName"`
	Email           string    `db:"email" json:"email" validate:"required,email"`
	Password        string    `db:"password" json:"password" validate:"password"`
	FirstName       string    `db:"firstname" json:"firstName"`
	LastName        string    `db:"lastname" json:"lastName"`
	BirthDate       string    `db:"birthdate" json:"birthDate"`
	Avatar          string    `db:"avatar" json:"avatar"`
	Role            Role      `db:"role" json:"role"`
	ActivateLink    string    `db:"activatelink" json:"activateLink"`
	IsActivate      bool      `db:"isactivate" json:"isActivate"`
	IsBanned        bool      `db:"isbanned" json:"isBanned"`
	ChannelsCount   int       `db:"channelscount" json:"channelsCount"`
	CreateTimestamp time.Time `db:"createtimestamp" json:"createTimestamp"`
}

type FindUsers struct {
	Id              int       `db:"id" json:"id"`
	NickName        string    `db:"nickname" json:"nickName"`
	Email           string    `db:"email" json:"email"`
	FirstName       string    `db:"firstname" json:"firstName"`
	LastName        string    `db:"lastname" json:"lastName"`
	BirthDate       string    `db:"birthdate" json:"birthDate"`
	Role            Role      `db:"role" json:"role"`
	IsBanned        bool      `db:"isbanned" json:"isBanned"`
	ChannelsCount   int       `db:"channelscount" json:"channelsCount"`
	CreateTimestamp time.Time `db:"createtimestamp" json:"createTimestamp"`
}

type ChangePasswordRequest struct {
	Id          int    `json:"id" validate:"required,numeric,min=1"`
	OldPassword string `json:"oldPassword" validate:"required"`
	NewPassword string `json:"NewPassword" validate:"required"`
}

type UserCreateResponse struct {
	*TokenResponse
	UserId int
}

type UserResponse struct {
	*TokenResponse
	*Users
}
