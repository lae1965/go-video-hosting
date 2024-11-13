package service

import (
	"fmt"
	"go-video-hosting/internal/database"
	"go-video-hosting/internal/errors"
	"go-video-hosting/internal/model"
)

type ChannelService struct {
	dbChannel         database.Channel
	dbUser            database.Users
	createTransaction CallbackFunc
}

func NewChannelService(dbChannel database.Channel, dbUser database.Users, CreateTransaction CallbackFunc) *ChannelService {
	return &ChannelService{
		dbChannel:         dbChannel,
		dbUser:            dbUser,
		createTransaction: CreateTransaction,
	}
}

func (channelService *ChannelService) CreateChannel(userId int, title string, description string) (int, *errors.AppError) {
	var err *errors.AppError
	transaction, errTr := channelService.createTransaction()
	if errTr != nil {
		return 0, errors.New(errors.UnknownError, fmt.Sprintf("failed create transaction: %s", errTr.Error()))
	}

	defer func() {
		if err != nil {
			transaction.Rollback()
		} else {
			transaction.Commit()
		}
	}()

	channelId, err := channelService.dbChannel.CreateChannel(transaction, userId, title, description) // Создаем канал
	if err != nil {
		err.Message = fmt.Sprintf("wrong creating channel: %s", err.Message)
		return 0, err
	}

	if err = channelService.dbUser.ChangeChannelsCountOfUser(transaction, userId, true); err != nil { // Увеличиваем количество каналов в таблице users
		err.Message = fmt.Sprintf("wrong incrementing channels count: %s", err.Message)
		return 0, err
	}

	return channelId, nil
}

func (channelService *ChannelService) UpdateChannel(userId int, channelId int, data map[string]string) *errors.AppError {
	return channelService.dbChannel.UpdateChannel(userId, channelId, data)
}

func (channelService *ChannelService) DeleteChannel(channelId int) *errors.AppError {
	var err *errors.AppError
	transaction, errTr := channelService.createTransaction()
	if errTr != nil {
		return errors.New(errors.UnknownError, fmt.Sprintf("failed create transaction: %s", errTr.Error()))
	}

	defer func() {
		if err != nil {
			transaction.Rollback()
		} else {
			transaction.Commit()
		}
	}()

	// TODO - удалить все видео канала с gRPC - сервера
	userId, err := channelService.dbChannel.DeleteChannel(transaction, channelId) // Удаляем канал
	if err != nil {
		err.Message = fmt.Sprintf("wrong deleting channel: %s", err.Message)
		return err
	}

	if err = channelService.dbUser.ChangeChannelsCountOfUser(transaction, userId, false); err != nil { // Уменьшаем количество каналов в таблице users
		err.Message = fmt.Sprintf("wrong incrementing channels count: %s", err.Message)
		return err
	}

	return nil
}

func (channelService *ChannelService) ToggleSubscribe(userId, channelId int) (*model.SubscribeRespose, *errors.AppError) {
	var err *errors.AppError
	transaction, errTr := channelService.createTransaction()
	if errTr != nil {
		return nil, errors.New(errors.UnknownError, fmt.Sprintf("failed create transaction: %s", errTr.Error()))
	}

	defer func() {
		if err != nil {
			transaction.Rollback()
		} else {
			transaction.Commit()
		}
	}()

	isSubscribe, err := channelService.dbChannel.ToggleSubscribe(transaction, userId, channelId)
	if err != nil {
		err.Message = fmt.Sprintf("wrong toggleSubscribing: %s", err.Message)
		return nil, err
	}

	subscribersCount, err := channelService.dbChannel.ChangeSubscribersCount(transaction, channelId, !isSubscribe)
	if err != nil {
		err.Message = fmt.Sprintf("wrong changing subscribesCount: %s", err.Message)
		return nil, err
	}

	return &model.SubscribeRespose{
		IsSubscribe:      isSubscribe,
		SubscribersCount: subscribersCount,
	}, nil
}

func (channelService *ChannelService) GetChannelById(userId, channelId int) (*model.GetChannelResponse, *errors.AppError) {
	channel, err := channelService.dbChannel.GetChannelById(channelId)
	if err != nil {
		err.Message = fmt.Sprintf("wrong getting channel: %s", err.Message)
		return nil, err
	}

	isSubscribe, err := channelService.dbChannel.IsSubscribe(userId, channelId)
	if err != nil {
		err.Message = fmt.Sprintf("channel subscription request error: %s", err.Message)
		return nil, err
	}

	return &model.GetChannelResponse{
		GetChannelFromDB: channel,
		IsSubscribe:      isSubscribe,
		IdList:           fmt.Sprintf("%d_%d", channel.UserId, channel.Id),
	}, nil
}

func (channelService *ChannelService) GetAllChannelsOfUser(userId int) ([]*model.GetAllChannelsResponse, *errors.AppError) {
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

func (channelService *ChannelService) GetAllIdListOfUser(userId int) ([]string, *errors.AppError) {
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
