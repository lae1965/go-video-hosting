package model

import "time"

type Video struct {
	Id              int       `json:"id"`
	Title           string    `json:"title"`
	HashName        string    `json:"hashName"`
	Category        string    `json:"category"`
	Description     string    `json:"description"`
	Thumbnail       string    `json:"thumbnail"`
	Duration        string    `json:"duration"`
	ViewsCount      int       `json:"viewsCount"`
	LikesCount      int       `json:"likesCount"`
	DislikesCount   int       `json:"dislikesCount"`
	IdList          string    `json:"idList"`
	CreateTimestamp time.Time `json:"createTimestamp"`
}

type VideoHistoty struct {
	Id              int       `json:"id"`
	UserId          int       `json:"userId"`
	VideoId         int       `json:"videoId"`
	UpdateTimestamp time.Time `json:"updateTimestamp"`
}

type CreateVideo struct {
	Title       string `json:"title" validate:"required"`
	IdList      string `json:"idList" validate:"required, id_list_len_3"`
	Thumbnail   string `json:"thumbnail" validate:"required"`
	Duration    string `json:"duration" validate:"required"`
	HashName    string `json:"hashName"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

type UpdateVideo struct {
	IdList      string `json:"idList" validate:"required, id_list_len_4"`
	Title       string `json:"title"`
	Category    string `json:"category"`
	Description string `json:"description"`
}
