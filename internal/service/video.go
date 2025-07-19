package service

import (
	"context"
	"fmt"
	grpcclient "go-video-hosting/gRPC/client"
	"go-video-hosting/internal/database"
	"go-video-hosting/internal/errors"
	"go-video-hosting/internal/model"
	"mime/multipart"
	"strconv"
	"strings"

	"cnb.cool/ordermap/ordermap"
	"github.com/sirupsen/logrus"
)

type VideoService struct {
	dbVideo           database.Video
	dbVideoHistory    database.VideoHistoty
	createTransaction CallbackFunc
	grpcClient        grpcclient.FilesGRPCClient
}

func NewVideoService(dbVideo database.Video, dbVideoHistory database.VideoHistoty, createTransaction CallbackFunc, grpcClient grpcclient.FilesGRPCClient) *VideoService {
	return &VideoService{
		dbVideo:           dbVideo,
		dbVideoHistory:    dbVideoHistory,
		createTransaction: createTransaction,
		grpcClient:        grpcClient,
	}
}

func (s *VideoService) SaveNewVideo(video *model.CreateVideo, fileHeader *multipart.FileHeader) (int, *errors.AppError) {
	var errDb *errors.AppError

	transaction, err := s.createTransaction()
	if err != nil {
		return 0, errors.New(errors.UnknownError, fmt.Sprintf("failed create transaction: %s", err.Error()))
	}

	idListArr := strings.Split(video.IdList, "_")
	channelId, _ := strconv.ParseInt(idListArr[1], 10, 0)

	isUnique, err := s.dbVideo.IsTitleUniqueForChannel(int(channelId), video.Title)
	if err != nil {
		return 0, errors.New(errors.UnknownError, err.Error())
	}
	if !isUnique {
		return 0, errors.New(errors.NotUnique, "channel's video title name must be unique")
	}

	videoHashName, err := s.grpcClient.SendToGRPCServer(context.Background(), fileHeader)
	if err != nil {
		return 0, errors.New(errors.UnknownError, fmt.Sprintf("can't save file to gRPC-server: %s", err.Error()))
	}

	video.HashName = videoHashName
	videoId, errDb := s.dbVideo.CreateNewVideo(transaction, video)
	if errDb != nil {
		errDb.Message = fmt.Sprintf("wrong saving video to db: %s", errDb.Message)
		return 0, errDb
	}

	idList := fmt.Sprintf("%s_%d", video.IdList, videoId)
	om := ordermap.New()
	om.Store("idList", idList)
	if errDb = s.dbVideo.UpdateVideoInfo(transaction, videoId, om); errDb != nil {
		errDb.Message = fmt.Sprintf("wrong saving idList to db: %s", errDb.Message)
		return 0, errDb
	}

	defer func() {
		if errDb != nil {
			transaction.Rollback()
			if videoHashName != "" {
				s.grpcClient.DeleteFromGRPCServer(context.Background(), videoHashName)
			}
		} else {
			transaction.Commit()
		}
	}()

	return videoId, nil
}

func (s *VideoService) GetVideoHashNameById(id int) (string, *errors.AppError) {
	return s.dbVideo.GetVideoHashNameById(id)
}

func (s *VideoService) DownloadVideo(hashName string, userId, videoId int, start, end int64, sendChunk func(int64, string, []byte) error) *errors.AppError {
	if err := s.grpcClient.GetFromGRPCServer(context.Background(), hashName, start, end, sendChunk); err != nil {
		return errors.New(errors.UnknownError, fmt.Sprintf("can't get video: %s", err.Error()))
	}

	if start == 0 {
		if err := s.dbVideo.IncrementViewsCount(videoId); err != nil {
			logrus.Errorf("error incrementing views count: %s", err.Error())
		}
		if err := s.dbVideoHistory.CreateVideoHistory(userId, videoId); err != nil {
			logrus.Errorf("error adding to browsing history: %s", err.Error())
		}
	}

	return nil
}

func (s *VideoService) UpdateVideoInfo(id int, data *ordermap.OrderMap) *errors.AppError {
	return s.dbVideo.UpdateVideoInfo(nil, id, data)
}
