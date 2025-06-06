package model

import "time"

type Playlist struct {
	Id              int       `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	IdList          string    `json:"idList"`
	CreateTimestamp time.Time `json:"createTimestamp"`
}

type CreatePlaylist struct {
	IdList      string `json:"idList" validate:"requied, id_list_len_2"`
	Title       string `json:"title" validate:"requied"`
	Description string `json:"description"`
}

type UpdatePlaylist struct {
	IdList         string `json:"idList" validate:"requied, id_list_len_3"`
	UpdatingObject struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
	}
}
