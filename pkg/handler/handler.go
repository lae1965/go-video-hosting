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
		userPublic := api.Group("/user")
		{
			userPublic.POST("/registration", handler.registration)
			userPublic.POST("/login", handler.login)
			userPublic.GET("/refresh", handler.refresh)
			userPublic.DELETE("/:id", handler.deleteUser) // TODO
		}

		user := api.Group("/user", handler.AuthMiddleware)
		{
			user.POST("/logout", handler.logout)
			user.PATCH("/edit/:id", handler.editUser)
			user.GET("/activate/:link", handler.activate) // TODO
			user.GET("/find_min/:id", handler.findMin)    // TODO
			user.GET("/find/:id", handler.find)           // TODO
			user.GET("/", handler.findAll)                // TODO
			user.POST("/avatar/:id", handler.saveAvatar)
			user.GET("/avatar/:id", handler.getAvatar)
			user.DELETE("/avatar/:id", handler.deleteAvatar)
			user.GET("/check", handler.checkPassword)              // TODO
			user.PATCH("/change_password", handler.changePassword) // TODO
		}
	}

	return router
}
