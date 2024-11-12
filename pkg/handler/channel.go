package handler

import (
	"go-video-hosting/internal/errors"
	"go-video-hosting/pkg/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) createChannel(ctx *gin.Context) {
	var input *model.CreateChannel

	if err := ctx.BindJSON(&input); err != nil {
		errors.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.validators.Validate.Struct(input); err != nil {
		errors.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userIdStr := strings.Split(input.IdList, "_")[0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 0)

	channelId, err := handler.services.Channel.CreateChannel(int(userId), input.Title, input.Description)
	if err != nil {
		errors.NewErrorResponse(ctx, err.Code, err.Message)
		return
	}

	ctx.JSON(http.StatusCreated, channelId)
}

func (handler *Handler) editChannel(ctx *gin.Context) {
	var input *model.UpdateChannel

	if err := ctx.BindJSON(&input); err != nil {
		errors.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.validators.Validate.Struct(input); err != nil {
		errors.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	idListArray := strings.Split(input.IdList, "_")
	userId, _ := strconv.ParseInt(idListArray[0], 10, 0)
	channelId, _ := strconv.ParseInt(idListArray[1], 10, 0)

	if err := handler.services.Channel.UpdateChannel(int(userId), int(channelId), map[string]string{
		"title":       input.UpdatingObject.Title,
		"description": input.UpdatingObject.Description,
	}); err != nil {
		errors.NewErrorResponse(ctx, err.Code, err.Message)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "updating success"})
}

func (handler *Handler) subscribe(ctx *gin.Context) {
	var input *model.SubscribeRequest

	if err := ctx.BindJSON(&input); err != nil {
		errors.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.validators.Validate.Struct(input); err != nil {
		errors.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response, err := handler.services.ToggleSubscribe(input.UserId, input.ChannelId)
	if err != nil {
		errors.NewErrorResponse(ctx, err.Code, err.Message)
	}

	ctx.JSON(http.StatusOK, response)
}

func (handler *Handler) removeChannel(ctx *gin.Context) {
	id, err := handler.GetIdFromQuery("id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		errors.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.services.DeleteChannel(id); err != nil {
		errors.NewErrorResponse(ctx, err.Code, err.Message)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (handler *Handler) getOneChannel(ctx *gin.Context) {
	getQuery := func(key string) string {
		return ctx.Query(key)
	}

	userId, err := handler.GetIdFromQuery("user_id", 0, getQuery)
	if err != nil {
		errors.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	channelId, err := handler.GetIdFromQuery("channel_id", 1, getQuery)
	if err != nil {
		errors.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	channel, errRes := handler.services.Channel.GetChannelById(userId, channelId)
	if errRes != nil {
		errors.NewErrorResponse(ctx, errRes.Code, errRes.Message)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"idList":           channel.IdList,
		"title":            channel.Title,
		"description":      channel.Description,
		"subscribersCount": channel.SubscribersCount,
		"createdTimestamp": channel.CreateTimestamp,
		"isSubscribed":     channel.IsSubscribe,
	})
}

func (handler *Handler) getAllChannelsOfUser(ctx *gin.Context) {
	userId, err := handler.GetIdFromQuery("user_id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		errors.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	channels, errRes := handler.services.Channel.GetAllChannelsOfUser(userId)
	if errRes != nil {
		errors.NewErrorResponse(ctx, errRes.Code, errRes.Message)
		return
	}

	ctx.JSON(http.StatusOK, channels)
}

func (handler *Handler) getSubscribersList(ctx *gin.Context) {
	userId, err := handler.GetIdFromQuery("user_id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		errors.NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	idLists, errRes := handler.services.Channel.GetAllIdListOfUser(userId)
	if errRes != nil {
		errors.NewErrorResponse(ctx, errRes.Code, errRes.Message)
		return
	}

	ctx.JSON(http.StatusOK, idLists)
}
