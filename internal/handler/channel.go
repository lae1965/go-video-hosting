package handler

import (
	"go-video-hosting/internal/errors"
	"go-video-hosting/internal/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createChannel(ctx *gin.Context) {
	var input *model.CreateChannel

	if err := ctx.BindJSON(&input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validators.Validate.Struct(input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userIdStr := strings.Split(input.IdList, "_")[0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 0)

	channelId, err := h.services.Channel.CreateChannel(int(userId), input.Title, input.Description)
	if err != nil {
		ErrorResponse(ctx, h.ErrorType2RequestStatus(err.Type), err.Message)
		return
	}

	ctx.JSON(http.StatusCreated, channelId)
}

func (h *Handler) editChannel(ctx *gin.Context) {
	var input *model.UpdateChannel

	if err := ctx.BindJSON(&input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validators.Validate.Struct(input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	idListArray := strings.Split(input.IdList, "_")
	userId, _ := strconv.ParseInt(idListArray[0], 10, 0)
	channelId, _ := strconv.ParseInt(idListArray[1], 10, 0)

	if err := h.services.Channel.UpdateChannel(int(userId), int(channelId), map[string]string{
		"title":       input.UpdatingObject.Title,
		"description": input.UpdatingObject.Description,
	}); err != nil {
		ErrorResponse(ctx, h.ErrorType2RequestStatus(err.Type), err.Message)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "updating success"})
}

func (h *Handler) subscribe(ctx *gin.Context) {
	var input *model.SubscribeRequest

	if err := ctx.BindJSON(&input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validators.Validate.Struct(input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response, err := h.services.ToggleSubscribe(input.UserId, input.ChannelId)
	if err != nil {
		ErrorResponse(ctx, h.ErrorType2RequestStatus(err.Type), err.Message)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) removeChannel(ctx *gin.Context) {
	id, err := h.GetIdFromQuery("id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.DeleteChannel(id); err != nil {
		ErrorResponse(ctx, h.ErrorType2RequestStatus(err.Type), err.Message)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *Handler) getOneChannel(ctx *gin.Context) {
	getQuery := func(key string) string {
		return ctx.Query(key)
	}

	userId, err := h.GetIdFromQuery("user_id", 0, getQuery)
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	channelId, err := h.GetIdFromQuery("channel_id", 1, getQuery)
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	channel, appErr := h.services.Channel.GetChannelById(userId, channelId)
	if appErr != nil {
		ErrorResponse(ctx, h.ErrorType2RequestStatus(appErr.Type), appErr.Message)
		return
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

func (h *Handler) getAllChannelsOfUser(ctx *gin.Context) {
	userId, err := h.GetIdFromQuery("user_id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	channels, appErr := h.services.Channel.GetAllChannelsOfUser(userId)
	if appErr != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, appErr.Message)
		return
	}

	if len(channels) == 0 {
		ctx.Status(http.StatusNoContent)
	} else {
		ctx.JSON(http.StatusOK, channels)
	}
}

func (h *Handler) getSubscribersList(ctx *gin.Context) {
	userId, err := h.GetIdFromQuery("user_id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	idLists, appErr := h.services.Channel.GetAllIdListOfUser(userId)
	if appErr != nil {
		var code int
		switch appErr.Type {
		case errors.EmptyField:
			code = http.StatusNoContent
		default:
			code = http.StatusInternalServerError
		}

		ErrorResponse(ctx, code, appErr.Message)
		return
	}

	ctx.JSON(http.StatusOK, idLists)
}
