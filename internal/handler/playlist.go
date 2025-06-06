package handler

import (
	"go-video-hosting/internal/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createPlaylist(ctx *gin.Context) {
	var input *model.CreatePlaylist

	if err := ctx.BindJSON(&input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	playlistId, err := h.services.Playlist.CreatePlaylist(input.IdList, input.Title, input.Description)
	if err != nil {
		ErrorResponse(ctx, h.ErrorType2RequestStatus(err.Type), err.Message)
		return
	}

	ctx.JSON(http.StatusCreated, playlistId)
}

func (h *Handler) editPlaylist(ctx *gin.Context) {
	var input *model.UpdatePlaylist

	if err := ctx.BindJSON(&input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validators.Validate.Struct(input); err != nil {
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

	if err := h.services.Playlist.UpdatePlaylist(input.IdList, data); err != nil {
		ErrorResponse(ctx, h.ErrorType2RequestStatus(err.Type), err.Message)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "updating success"})
}

func (h *Handler) removePlaylist(ctx *gin.Context) {
	id, err := h.GetIdFromQuery("id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.DeletePlaylist(id); err != nil {
		ErrorResponse(ctx, h.ErrorType2RequestStatus(err.Type), err.Message)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (h *Handler) getOnePlaylist(ctx *gin.Context) {
	playlistId, err := h.GetIdFromQuery("id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	playlist, appErr := h.services.Playlist.GetPlaylistById(playlistId)
	if appErr != nil {
		ErrorResponse(ctx, h.ErrorType2RequestStatus(appErr.Type), appErr.Message)
		return
	}

	ctx.JSON(http.StatusOK, playlist)
}

func (h *Handler) getAllPlaylistsOfChannel(ctx *gin.Context) {
	channelId, err := h.GetIdFromQuery("channel_id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	playlists, appErr := h.services.Playlist.GetAllPlaylistsOfChannel(channelId)
	if appErr != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, appErr.Message)
		return
	}

	if len(playlists) == 0 {
		ctx.Status(http.StatusNoContent)
	} else {
		ctx.JSON(http.StatusOK, playlists)
	}
}
