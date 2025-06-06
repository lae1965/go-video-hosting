package handler

import (
	"go-video-hosting/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) createPlaylist(ctx *gin.Context) {
	var input *model.CreatePlaylist

	if err := ctx.BindJSON(&input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	playlistId, err := handler.services.Playlist.CreatePlaylist(input.IdList, input.Title, input.Description)
	if err != nil {
		ErrorResponse(ctx, handler.ErrorType2RequestStatus(err.Type), err.Message)
		return
	}

	ctx.JSON(http.StatusCreated, playlistId)
}

func (handler *Handler) editPlaylist(ctx *gin.Context) {
	var input *model.UpdatePlaylist

	if err := ctx.BindJSON(&input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.validators.Validate.Struct(input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	data := make(map[string]string)
	if input.UpdatingObject.Title != nil {
		if *input.UpdatingObject.Title == "" {
			ErrorResponse(ctx, http.StatusBadRequest, "Title must not be empty")
			return
		}
		data["title"] = *input.UpdatingObject.Title
	}
	if input.UpdatingObject.Description != nil {
		data["description"] = *input.UpdatingObject.Description
	}

	if err := handler.services.Playlist.UpdatePlaylist(input.IdList, data); err != nil {
		ErrorResponse(ctx, handler.ErrorType2RequestStatus(err.Type), err.Message)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "updating success"})
}
