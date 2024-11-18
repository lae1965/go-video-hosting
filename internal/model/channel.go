package model

import (
	_ "github.com/go-playground/validator/v10"

	"time"
)

type Channel struct {
	Id               int       `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	SubscribersCount int       `json:"subscribersCount"`
	CreateTimestamp  time.Time `json:"createTimestamp"`
}

type CreateChannel struct {
	IdList      string `json:"idList" validate:"required,channelIdList"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description"`
}

type UpdateChannel struct {
	IdList         string `json:"idList" validate:"required,channelIdList"`
	UpdatingObject struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
}

type SubscribeRequest struct {
	UserId    int `json:"userId" validate:"required,min=1"`
	ChannelId int `json:"id" validate:"required,min=1"`
}

type SubscribeRespose struct {
	IsSubscribe      bool `json:"isSubscribe"`
	SubscribersCount int  `json:"subscribersCount"`
}

type GetChannelFromDB struct {
	*Channel
	UserId int `json:"userId"`
}

type GetChannelResponse struct {
	*GetChannelFromDB
	IdList      string `json:"idList"`
	IsSubscribe bool   `json:"isSubscribe"`
}

type GetAllChannelsResponse struct {
	IdList           string    `json:"idList"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	SubscribersCount int       `json:"subscribersCount"`
	CreatedTimestamp time.Time `json:"createdTimestamp"`
}
