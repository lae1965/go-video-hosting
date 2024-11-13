package handler

import (
	"fmt"
	"go-video-hosting/internal/service"
	"go-video-hosting/internal/validator"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services   *service.Service
	validators *validator.Validator
}

func NewHandler(services *service.Service, validators *validator.Validator) *Handler {
	return &Handler{services: services, validators: validators}
}

func (handler *Handler) GetIdFromQuery(key string, min int, getKey func(string) string) (int, error) {
	idStr := getKey(key)
	if err := handler.validators.Validate.Var(idStr, fmt.Sprintf("required,numeric,min=%d", min)); err != nil {
		return 0, err
	}

	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (handler *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("api")
	{
		userPublic := api.Group("/user")
		{
			userPublic.POST("/registration", handler.registration)
			userPublic.POST("/login", handler.login)
			userPublic.GET("/refresh", handler.refresh)
			userPublic.GET("/activate/:link", handler.activate)
		}

		// user := api.Group("/user", handler.AuthMiddleware) //! for testing
		user := api.Group("/user") //! for testing
		{
			user.POST("/logout", handler.logout)
			user.PATCH("/edit/:id", handler.editUser)
			user.DELETE("/:id", handler.deleteUser)
			user.GET("/find_min/:id", handler.findMin)
			user.GET("/find/:id", handler.find)
			user.GET("", handler.findAll)
			user.POST("/avatar/:id", handler.saveAvatar)
			user.GET("/avatar/:id", handler.getAvatar)
			user.DELETE("/avatar/:id", handler.deleteAvatar)
			user.GET("/check", handler.checkUnique)
			user.PATCH("/change_password", handler.changePassword)
		}

		// channel := api.Group("/channel", handler.AuthMiddleware) //! for testing
		channel := api.Group("/channel") //! for testing
		{
			channel.POST("/create", handler.createChannel)
			channel.PATCH("/edit", handler.editChannel)
			channel.PATCH("/subscribe", handler.subscribe)
			channel.DELETE("/:id", handler.removeChannel)
			channel.GET("/get_one", handler.getOneChannel)
			channel.GET("/get_all/:user_id", handler.getAllChannelsOfUser)
			channel.GET("subscribes_list/:user_id", handler.getSubscribersList)
		}
	}

	return router
}
