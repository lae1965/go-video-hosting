package service

import (
	"fmt"
	"go-video-hosting/internal/database"
	"go-video-hosting/internal/errors"
	"go-video-hosting/internal/model"
)

type ChannelService struct {
	dbChannel database.Channel
	dbUser    database.Users
}

func NewChannelService(dbChannel database.Channel, dbUser database.Users) *ChannelService {
	return &ChannelService{
		dbChannel: dbChannel,
		dbUser:    dbUser,
	}
}

func (channelService *ChannelService) CreateChannel(userId int, title string, description string) (int, *errors.ErrorRes) {
	channelId, err := channelService.dbChannel.CreateChannel(userId, title, description) // Создаем канал
	if err != nil {
		return 0, err
	}

	if err = channelService.dbUser.ChangeChannelsCountOfUser(userId, true); err != nil { // Увеличиваем количество каналов в таблице users
		return 0, err
	}

	return channelId, nil
}

func (channelService *ChannelService) UpdateChannel(userId int, channelId int, data map[string]string) *errors.ErrorRes {
	return channelService.dbChannel.UpdateChannel(userId, channelId, data)
}

func (channelService *ChannelService) DeleteChannel(channelId int) *errors.ErrorRes {
	// TODO - удалить все видео канала с gRPC - сервера
	userId, err := channelService.dbChannel.DeleteChannel(channelId) // Удаляем канал
	if err != nil {
		return err
	}

	if err = channelService.dbUser.ChangeChannelsCountOfUser(userId, false); err != nil { // Уменьшаем количество каналов в таблице users
		return err
	}

	return nil
}

func (channelService *ChannelService) ToggleSubscribe(userId, channelId int) (*model.SubscribeRespose, *errors.ErrorRes) {
	isSubscribe, err := channelService.dbChannel.ToggleSubscribe(userId, channelId)
	if err != nil {
		return nil, err
	}

	subscribersCount, err := channelService.dbChannel.ChangeSubscribersCount(channelId, !isSubscribe)
	if err != nil {
		return nil, err
	}

	return &model.SubscribeRespose{
		IsSubscribe:      isSubscribe,
		SubscribersCount: subscribersCount,
	}, nil
}

func (channelService *ChannelService) GetChannelById(userId, channelId int) (*model.GetChannelResponse, *errors.ErrorRes) {
	channel, err := channelService.dbChannel.GetChannelById(channelId)
	if err != nil {
		return nil, err
	}

	isSubscribe, err := channelService.dbChannel.IsSubscribe(userId, channelId)
	if err != nil {
		return nil, err
	}

	return &model.GetChannelResponse{
		GetChannelFromDB: channel,
		IsSubscribe:      isSubscribe,
		IdList:           fmt.Sprintf("%d_%d", channel.UserId, channel.Id),
	}, nil
}

func (channelService *ChannelService) GetAllChannelsOfUser(userId int) ([]*model.GetAllChannelsResponse, *errors.ErrorRes) {
	channelsFromDb, err := channelService.dbChannel.GetAllChannelsOfUser(userId)
	if err != nil {
		return nil, err
	}

	var channelsResponse []*model.GetAllChannelsResponse

	for _, channel := range channelsFromDb {
		channelsResponse = append(channelsResponse, &model.GetAllChannelsResponse{
			IdList:           fmt.Sprintf("%d_%d", userId, channel.Id),
			Title:            channel.Title,
			Description:      channel.Description,
			SubscribersCount: channel.SubscribersCount,
			CreatedTimestamp: channel.CreateTimestamp,
		})
	}

	return channelsResponse, nil
}

func (channelService *ChannelService) GetAllIdListOfUser(userId int) ([]string, *errors.ErrorRes) {
	idsList, err := channelService.dbChannel.GetSubscribingChannelsOfUser(userId)
	if err != nil {
		return nil, err
	}

	var idLists []string
	for _, ids := range idsList {
		idLists = append(idLists, fmt.Sprintf("%d_%d", ids.UserId, ids.ChannelId))
	}

	return idLists, nil
}
