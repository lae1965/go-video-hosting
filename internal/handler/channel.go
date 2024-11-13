package handler

import (
	"go-video-hosting/internal/errors"
	"go-video-hosting/internal/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) createChannel(ctx *gin.Context) {
	var input *model.CreateChannel

	if err := ctx.BindJSON(&input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.validators.Validate.Struct(input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userIdStr := strings.Split(input.IdList, "_")[0]
	userId, _ := strconv.ParseInt(userIdStr, 10, 0)

	channelId, err := handler.services.Channel.CreateChannel(int(userId), input.Title, input.Description)
	if err != nil {
		var code int
		switch err.Type {
		case errors.NotFound:
			code = http.StatusBadRequest
		case errors.NotUnique:
			code = http.StatusConflict
		default:
			code = http.StatusInternalServerError
		}

		ErrorResponse(ctx, code, err.Message)
		return
	}

	ctx.JSON(http.StatusCreated, channelId)
}

func (handler *Handler) editChannel(ctx *gin.Context) {
	var input *model.UpdateChannel

	if err := ctx.BindJSON(&input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.validators.Validate.Struct(input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	idListArray := strings.Split(input.IdList, "_")
	userId, _ := strconv.ParseInt(idListArray[0], 10, 0)
	channelId, _ := strconv.ParseInt(idListArray[1], 10, 0)

	if err := handler.services.Channel.UpdateChannel(int(userId), int(channelId), map[string]string{
		"title":       input.UpdatingObject.Title,
		"description": input.UpdatingObject.Description,
	}); err != nil {
		var code int
		switch err.Type {
		case errors.NotFound:
			code = http.StatusBadRequest
		case errors.NotUnique:
			code = http.StatusConflict
		default:
			code = http.StatusInternalServerError
		}

		ErrorResponse(ctx, code, err.Message)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "updating success"})
}

func (handler *Handler) subscribe(ctx *gin.Context) {
	var input *model.SubscribeRequest

	if err := ctx.BindJSON(&input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.validators.Validate.Struct(input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	response, err := handler.services.ToggleSubscribe(input.UserId, input.ChannelId)
	if err != nil {
		var code int
		switch err.Type {
		case errors.NotFound:
			code = http.StatusNotFound
		default:
			code = http.StatusInternalServerError
		}

		ErrorResponse(ctx, code, err.Message)
	}

	ctx.JSON(http.StatusOK, response)
}

func (handler *Handler) removeChannel(ctx *gin.Context) {
	id, err := handler.GetIdFromQuery("id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.services.DeleteChannel(id); err != nil {
		var code int
		switch err.Type {
		case errors.NotFound:
			code = http.StatusNotFound
		default:
			code = http.StatusInternalServerError
		}

		ErrorResponse(ctx, code, err.Message)
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
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	channelId, err := handler.GetIdFromQuery("channel_id", 1, getQuery)
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	channel, appErr := handler.services.Channel.GetChannelById(userId, channelId)
	if appErr != nil {
		var code int
		switch appErr.Type {
		case errors.NotFound:
			code = http.StatusNotFound
		default:
			code = http.StatusInternalServerError
		}

		ErrorResponse(ctx, code, appErr.Message)
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
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	channels, appErr := handler.services.Channel.GetAllChannelsOfUser(userId)
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

	ctx.JSON(http.StatusOK, channels)
}

func (handler *Handler) getSubscribersList(ctx *gin.Context) {
	userId, err := handler.GetIdFromQuery("user_id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	idLists, appErr := handler.services.Channel.GetAllIdListOfUser(userId)
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
