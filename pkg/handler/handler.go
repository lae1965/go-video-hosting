package handler

import (
	"go-video-hosting/pkg/service"
	"go-video-hosting/pkg/validator"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services   *service.Service
	validators *validator.Validator
}

func NewHandler(services *service.Service, validators *validator.Validator) *Handler {
	return &Handler{services: services, validators: validators}
}

func (handler *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("api")
	{
		user := api.Group("/user")
		{
			user.POST("/registration", handler.registration) //! Done
			user.POST("/login", handler.login)               //! Done
			user.POST("/logout", handler.logout)             //! Done
			user.PATCH("/edit/:id", handler.editUser)
			user.DELETE("/:id", handler.deleteUser)
			user.GET("/refresh", handler.refresh) //! Done
			user.GET("/activate/:link", handler.activate)
			user.GET("/find_min/:id", handler.findMin)
			user.GET("/find/:id", handler.find)
			user.GET("/", handler.findAll)
			user.POST("/avatar/:id", handler.saveAvatar)     //! Done
			user.GET("/avatar/:id", handler.getAvatar)       //! Done
			user.DELETE("/avatar/:id", handler.deleteAvatar) //! Done
			user.GET("/check", handler.checkPassword)
			user.PATCH("/change_password", handler.changePassword)
		}
	}

	return router
}
