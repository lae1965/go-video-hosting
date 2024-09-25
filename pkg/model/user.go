package model

import "time"

type Role int

const (
	UserRole Role = iota
	AdminRole
	SuperAdminRole
)

type Users struct {
	Id              int       `json:"-"`
	NickName        string    `json:"nickName"`
	Email           string    `json:"email"`
	Password        string    `json:"password"`
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
