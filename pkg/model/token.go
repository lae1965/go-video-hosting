package model

type Token struct {
	Id     int    `json:"id"`
	Token  string `json:"token"`
	UserId int    `json:"userId"`
}
