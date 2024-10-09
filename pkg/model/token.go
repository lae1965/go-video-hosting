package model

type Token struct {
	Id     int    `json:"id"`
	Token  string `json:"token"`
	UserId int    `json:"userId"`
}

type TokenResponse struct {
	AccessToken    string
	RefreshToken   string
	RefreshTokenId int
}
