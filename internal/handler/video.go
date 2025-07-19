package handler

import (
	"fmt"
	"go-video-hosting/internal/model"
	"net/http"
	"strconv"
	"strings"

	"cnb.cool/ordermap/ordermap"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createVideo(ctx *gin.Context) {
	var input *model.CreateVideo

	if err := ctx.BindJSON(&input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	fileHeader, err := ctx.FormFile("files")
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.validators.Validate.Var(fileHeader, "videofile"); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("invalid file type: %s", err.Error()))
		return
	}

	videoId, appErr := h.services.Video.SaveNewVideo(input, fileHeader)
	if appErr != nil {
		ErrorResponse(ctx, h.ErrorType2RequestStatus(appErr.Type), appErr.Message)
		return
	}

	ctx.JSON(http.StatusCreated, videoId)
}

func (h *Handler) getVideoName(ctx *gin.Context) {
	id, err := h.GetIdFromQuery("id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	hashName, appErr := h.services.Video.GetVideoHashNameById(id)
	if appErr != nil {
		ErrorResponse(ctx, h.ErrorType2RequestStatus(appErr.Type), appErr.Message)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"hashName": hashName})
}

func (h *Handler) downloadVideo(ctx *gin.Context) {
	hashName := ctx.Query("hash_name")
	if err := h.validators.Validate.Var(hashName, "hash"); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("invalid file type: %s", err.Error()))
		return
	}

	userId, err := h.GetIdFromQuery("user_id", 1, func(key string) string {
		return ctx.Query(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	videoId, err := h.GetIdFromQuery("video_id", 1, func(key string) string {
		return ctx.Query(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	rangeHeader := ctx.GetHeader("Range")
	var start, end int64
	if rangeHeader == "" {
		start = 0
		end = -1
	} else {
		rangeWithoutBytes := strings.Replace(rangeHeader, "bytes=", "", -1)
		rangeArray := strings.Split(rangeWithoutBytes, "-")
		start, err = strconv.ParseInt(rangeArray[0], 10, 64)
		if err != nil {
			ErrorResponse(ctx, http.StatusBadRequest, err.Error())
			return
		}
		if len(rangeArray) > 1 && rangeArray[1] != "" {
			end, err = strconv.ParseInt(rangeArray[1], 10, 64)
			if err != nil {
				ErrorResponse(ctx, http.StatusBadRequest, err.Error())
				return
			}
		} else {
			end = -1
		}
	}

	isHeadersSet := false
	appErr := h.services.Video.DownloadVideo(hashName, userId, videoId, start, end, func(fileSize int64, mimeType string, chunk []byte) error {
		if !isHeadersSet {
			if rangeHeader == "" {
				ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
				ctx.Writer.WriteHeader(http.StatusOK)
			} else {
				if end == -1 {
					end = fileSize - 1
				}
				ctx.Writer.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
				ctx.Writer.Header().Set("Accept-Ranges", "bytes")
				ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", end-start+1))
				ctx.Writer.WriteHeader(http.StatusPartialContent)
			}
			ctx.Writer.Header().Set("Content-Type", mimeType)
			isHeadersSet = true
		}
		_, err := ctx.Writer.Write(chunk)

		return err
	})
	if appErr != nil {
		ErrorResponse(ctx, h.ErrorType2RequestStatus(appErr.Type), appErr.Message)
		return
	}
}

func (h *Handler) editVideoInfo(ctx *gin.Context) {
	var jsonObject, video map[string]interface{}
	if err := ctx.BindJSON(&jsonObject); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if updatingObject, exists := jsonObject["updatingObject"]; exists {
		var ok bool
		if video, ok = updatingObject.(map[string]interface{}); !ok {
			ErrorResponse(ctx, http.StatusBadRequest, "invalid updatingObject")
			return
		}
	} else {
		ErrorResponse(ctx, http.StatusBadRequest, `body must contain an object with the key "updatingObject". Data must be inside this object`)
		return
	}

	omVideo := ordermap.New()
	for key, value := range video {
		omVideo.Store(key, value)
	}

	idList := video["idList"].(string)
	idListArray := strings.Split(idList, "_")
	videoId, err := strconv.ParseInt(idListArray[4], 10, 64)
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if appErr := h.services.UpdateUser(int(videoId), omVideo); appErr != nil {
		ErrorResponse(ctx, h.ErrorType2RequestStatus(appErr.Type), appErr.Message)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "update success"})
}

func (h *Handler) deleteVideo(ctx *gin.Context)              {}
func (h *Handler) getVideoInfoById(ctx *gin.Context)         {}
func (h *Handler) getFavoriteVideo(ctx *gin.Context)         {}
func (h *Handler) getVideosOfChannel(ctx *gin.Context)       {}
func (h *Handler) getVideosOfPlaylist(ctx *gin.Context)      {}
func (h *Handler) getFrameshortbyVideoname(ctx *gin.Context) {}
func (h *Handler) toggleLike(ctx *gin.Context)               {}
func (h *Handler) toggleDislike(ctx *gin.Context)            {}
func (h *Handler) getVideosByFrame(ctx *gin.Context)         {}
func (h *Handler) getVideosViewHistory(ctx *gin.Context)     {}
func (h *Handler) getVListOfFavoriteVideos(ctx *gin.Context) {}
func (h *Handler) deleteViewFromHistory(ctx *gin.Context)    {}
func (h *Handler) deleteViewHistoryOfUser(ctx *gin.Context)  {}
