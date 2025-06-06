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

func (s *ChannelService) CreateChannel(userId int, title string, description string) (int, *errors.AppError) {
	var errDB *errors.AppError
	transaction, errTr := s.createTransaction()
	if errTr != nil {
		return 0, errors.New(errors.UnknownError, fmt.Sprintf("failed create transaction: %s", errTr.Error()))
	}

	defer func() {
		if errDB != nil {
			transaction.Rollback()
		} else {
			transaction.Commit()
		}
	}()

	isUserExist, err := s.dbChannel.IsUserExist(userId)
	if err != nil {
		return 0, errors.New(errors.UnknownError, err.Error())
	}
	if !isUserExist {
		return 0, errors.New(errors.NotFound, fmt.Sprintf("user with userid = %d not found", userId))
	}

	isUnique, err := s.dbChannel.IsTitlelUniqueForUser(userId, title)
	if err != nil {
		return 0, errors.New(errors.UnknownError, err.Error())
	}
	if !isUnique {
		return 0, errors.New(errors.NotUnique, "user's channel name must be unique")
	}

	channelId, errDB := s.dbChannel.CreateChannel(transaction, userId, title, description) // Создаем канал
	if errDB != nil {
		errDB.Message = fmt.Sprintf("wrong creating channel: %s", errDB.Message)
		return 0, errDB
	}

	if errDB = s.dbUser.ChangeChannelsCountOfUser(transaction, userId, true); errDB != nil { // Увеличиваем количество каналов в таблице users
		errDB.Message = fmt.Sprintf("wrong incrementing channels count: %s", errDB.Message)
		return 0, errDB
	}

	return channelId, nil
}

func (s *ChannelService) UpdateChannel(userId int, channelId int, data map[string]string) *errors.AppError {
	isUserExist, err := s.dbChannel.IsUserExist(userId)
	if err != nil {
		return errors.New(errors.UnknownError, err.Error())
	}
	if !isUserExist {
		return errors.New(errors.NotFound, fmt.Sprintf("user with userId = %d not exist", userId))
	}

	if title, exist := data["title"]; exist {
		isUnique, err := s.dbChannel.IsTitlelUniqueForUser(userId, title)
		if err != nil {
			return errors.New(errors.UnknownError, err.Error())
		}
		if !isUnique {
			return errors.New(errors.NotUnique, "user's channel name must be unique")
		}
	}

	return s.dbChannel.UpdateChannel(userId, channelId, data)
}

func (s *ChannelService) DeleteChannel(channelId int) *errors.AppError {
	var err *errors.AppError
	transaction, errTr := s.createTransaction()
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
	userId, err := s.dbChannel.DeleteChannel(transaction, channelId) // Удаляем канал
	if err != nil {
		err.Message = fmt.Sprintf("wrong deleting channel: %s", err.Message)
		return err
	}

	if err = s.dbUser.ChangeChannelsCountOfUser(transaction, userId, false); err != nil { // Уменьшаем количество каналов в таблице users
		err.Message = fmt.Sprintf("wrong incrementing channels count: %s", err.Message)
		return err
	}

	return nil
}

func (s *ChannelService) ToggleSubscribe(userId, channelId int) (*model.SubscribeRespose, *errors.AppError) {
	var err *errors.AppError
	transaction, errTr := s.createTransaction()
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

	isSubscribe, err := s.dbChannel.ToggleSubscribe(transaction, userId, channelId)
	if err != nil {
		err.Message = fmt.Sprintf("wrong toggleSubscribing: %s", err.Message)
		return nil, err
	}

	subscribersCount, err := s.dbChannel.ChangeSubscribersCount(transaction, channelId, !isSubscribe)
	if err != nil {
		err.Message = fmt.Sprintf("wrong changing subscribesCount: %s", err.Message)
		return nil, err
	}

	return &model.SubscribeRespose{
		IsSubscribe:      isSubscribe,
		SubscribersCount: subscribersCount,
	}, nil
}

func (s *ChannelService) GetChannelById(userId, channelId int) (*model.GetChannelResponse, *errors.AppError) {
	channel, err := s.dbChannel.GetChannelById(channelId)
	if err != nil {
		err.Message = fmt.Sprintf("wrong getting channel: %s", err.Message)
		return nil, err
	}

	isSubscribe, err := s.dbChannel.IsSubscribe(userId, channelId)
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

func (s *ChannelService) GetAllChannelsOfUser(userId int) ([]*model.GetAllChannelsResponse, *errors.AppError) {
	channelsFromDb, err := s.dbChannel.GetAllChannelsOfUser(userId)
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

func (s *ChannelService) GetAllIdListOfUser(userId int) ([]string, *errors.AppError) {
	idsList, err := s.dbChannel.GetSubscribingChannelsOfUser(userId)
	if err != nil {
		return nil, err
	}

	var idLists []string
	for _, ids := range idsList {
		idLists = append(idLists, fmt.Sprintf("%d_%d", ids.UserId, ids.ChannelId))
	}

	return idLists, nil
}
