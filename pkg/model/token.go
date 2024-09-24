package model

type Token struct {
	Id     int    `json:"-"`
	Token  string `json:"token"`
	UserId int    `json:"userId"`
}
