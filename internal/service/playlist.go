package service

import (
	"fmt"
	"go-video-hosting/internal/database"
	"go-video-hosting/internal/errors"
	"strconv"
	"strings"
)

type PlaylistService struct {
	dbPlaylist        database.Playlist
	createTransaction CallbackFunc
}

func NewPlaylistService(dbPlaylist database.Playlist, createTransaction CallbackFunc) *PlaylistService {
	return &PlaylistService{
		dbPlaylist:        dbPlaylist,
		createTransaction: createTransaction,
	}
}

func (s *PlaylistService) CreatePlaylist(idList string, title string, description string) (int, *errors.AppError) {
	var errDB *errors.AppError

	transaction, err := s.createTransaction()
	if err != nil {
		return 0, errors.New(errors.UnknownError, fmt.Sprintf("failed create transaction: %s", err.Error()))
	}

	defer func() {
		if errDB != nil {
			transaction.Rollback()
		} else {
			transaction.Commit()
		}
	}()

	idListArr := strings.Split(idList, "_")
	channelId, _ := strconv.ParseInt(idListArr[1], 10, 0)

	isChannelExist, err := s.dbPlaylist.IsChannelExist(int(channelId))
	if err != nil {
		return 0, errors.New(errors.UnknownError, err.Error())
	}
	if !isChannelExist {
		return 0, errors.New(errors.NotFound, fmt.Sprintf("channel with channelId = %d not exist", channelId))
	}

	isUnique, err := s.dbPlaylist.IsTitlelUniqueForChannel(int(channelId), title)
	if err != nil {
		return 0, errors.New(errors.UnknownError, err.Error())
	}
	if !isUnique {
		return 0, errors.New(errors.NotUnique, "channel's playlist name must be unique")
	}

	playlistId, errDB := s.dbPlaylist.CreatePlaylist(transaction, int(channelId), title, description)
	if errDB != nil {
		errDB.Message = fmt.Sprintf("wrong creating playlist: %s", errDB.Message)
		return 0, errDB
	}

	idList = fmt.Sprintf("%s_%d", idList, playlistId)
	errDB = s.dbPlaylist.UpdatePlaylist(transaction, int(channelId), playlistId, map[string]string{"idList": idList})

	return playlistId, nil
}

func (s *PlaylistService) UpdatePlaylist(idList string, data map[string]string) *errors.AppError {
	idListArr := strings.Split(idList, "_")
	channelId64, _ := strconv.ParseInt(idListArr[1], 10, 0)
	channelId := int(channelId64)
	playlistId64, _ := strconv.ParseInt(idListArr[2], 10, 0)
	playlistId := int(playlistId64)

	isChannelExist, err := s.dbPlaylist.IsChannelExist(channelId)
	if err != nil {
		return errors.New(errors.UnknownError, err.Error())
	}
	if !isChannelExist {
		return errors.New(errors.NotFound, fmt.Sprintf("channel with channelId = %d not exist", channelId))
	}

	if title, exist := data["title"]; exist {
		isUnique, err := s.dbPlaylist.IsTitlelUniqueForChannel(channelId, title)
		if err != nil {
			return errors.New(errors.UnknownError, err.Error())
		}
		if !isUnique {
			return errors.New(errors.NotUnique, "channel's playlist name must be unique")
		}
	}

	return s.dbPlaylist.UpdatePlaylist(nil, channelId, playlistId, data)
}
