package handler

import (
	"fmt"
	"go-video-hosting/internal/errors"
	"go-video-hosting/internal/service"
	"go-video-hosting/internal/validator"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services   *service.Service
	validators *validator.Validator
}

func New(services *service.Service, validators *validator.Validator) *Handler {
	return &Handler{services: services, validators: validators}
}

func (h *Handler) GetIdFromQuery(key string, min int, getKey func(string) string) (int, error) {
	idStr := getKey(key)
	if err := h.validators.Validate.Var(idStr, fmt.Sprintf("required,numeric,min=%d", min)); err != nil {
		return 0, err
	}

	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (h *Handler) ErrorType2RequestStatus(errType errors.ErrType) int {
	var code int
	switch errType {
	case errors.NotFound:
		code = http.StatusBadRequest
	case errors.NotEqual:
		code = http.StatusConflict
	case errors.NotUnique:
		code = http.StatusConflict
	case errors.Unauthorization:
		code = http.StatusUnauthorized
	case errors.EmptyField:
		code = http.StatusNoContent
	default:
		code = http.StatusInternalServerError
	}
	return code
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("api")
	{
		userPublic := api.Group("/user")
		{
			userPublic.POST("/registration", h.registration)
			userPublic.POST("/login", h.login)
			userPublic.GET("/refresh", h.refresh)
			userPublic.GET("/activate/:link", h.activate)
		}

		// user := api.Group("/user", h.AuthMiddleware) //! for testing
		user := api.Group("/user") //! for testing
		{
			user.POST("/logout", h.logout)
			user.PATCH("/edit/:id", h.editUser)
			user.DELETE("/:id", h.deleteUser)
			user.GET("/find_min/:id", h.getMin)
			user.GET("/find/:id", h.getById)
			user.GET("", h.getAll)
			user.POST("/avatar/:id", h.saveAvatar)
			user.GET("/avatar/:id", h.getAvatar)
			user.DELETE("/avatar/:id", h.deleteAvatar)
			user.GET("/check", h.checkUnique)
			user.PATCH("/change_password", h.changePassword)
		}

		// channel := api.Group("/channel", h.AuthMiddleware) //! for testing
		channel := api.Group("/channel") //! for testing
		{
			channel.POST("/create", h.createChannel)
			channel.PATCH("/edit", h.editChannel)
			channel.PATCH("/subscribe", h.subscribe)
			channel.DELETE("/:id", h.removeChannel)
			channel.GET("/get_one", h.getOneChannel)
			channel.GET("/get_all/:user_id", h.getAllChannelsOfUser)
			channel.GET("subscribes_list/:user_id", h.getSubscribersList)
		}

		// playlist := api.Group("/playlist", h.AuthMiddleware) //! for testing
		playlist := api.Group("/playlist") //! for testing
		{
			playlist.POST("/create", h.createPlaylist)
			playlist.PATCH("/edit", h.editPlaylist)
			playlist.DELETE("/:id", h.removePlaylist)
			playlist.GET("/get_one/:id", h.getOnePlaylist)
			playlist.GET("/get_all/:channel_id", h.getAllPlaylistsOfChannel)
		}

		// video := api.Group("/video", h.AuthMiddleware) //! for testing
		video := api.Group("/video") //! for testing
		{
			video.POST("/create", h.createVideo)
			video.GET("/get_name/:id", h.getVideoName)
			video.GET("/download", h.downloadVideo)
			video.PATCH("/edit", h.editVideoInfo)
			video.DELETE("/:id", h.deleteVideo)
			video.GET("/get_one", h.getVideoInfoById)
			video.GET("/get_favorite", h.getFavoriteVideo)
			video.GET("/channel/:id", h.getVideosOfChannel)
			video.GET("/get_all/:playlist_id", h.getVideosOfPlaylist)
			video.GET("/frameshort/:hash_name", h.getFrameshortbyVideoname)
			video.PATCH("/like", h.toggleLike)
			video.PATCH("/dislike", h.toggleDislike)
			video.GET("/query/:title", h.getVideosByFrame)
			video.GET("/history/:user_id", h.getVideosViewHistory)
			video.GET("/likes_list/:user_id", h.getVListOfFavoriteVideos)
			video.DELETE("/history/del_one", h.deleteViewFromHistory)
			video.DELETE("/history/del_all/:user_id", h.deleteViewHistoryOfUser)
		}
		return router
	}
}
