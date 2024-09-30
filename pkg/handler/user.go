package handler

import (
	"go-video-hosting/pkg/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (handler *Handler) registration(ctx *gin.Context) {
	var input model.Users

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.validators.Validate.Struct(input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := handler.services.CreateUser(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, map[string]interface{}{"id": id})
}

func (handler *Handler) login(ctx *gin.Context) {

}

func (handler *Handler) logout(ctx *gin.Context) {

}

func (handler *Handler) editUser(ctx *gin.Context) {

}

func (handler *Handler) deleteUser(ctx *gin.Context) {

}

func (handler *Handler) refresh(ctx *gin.Context) {

}

func (handler *Handler) activate(ctx *gin.Context) {

}

func (handler *Handler) findMin(ctx *gin.Context) {

}

func (handler *Handler) find(ctx *gin.Context) {

}

func (handler *Handler) findAll(ctx *gin.Context) {

}

func (handler *Handler) saveAvatar(ctx *gin.Context) {

}

func (handler *Handler) getAvatar(ctx *gin.Context) {

}

func (handler *Handler) deleteAvatar(ctx *gin.Context) {

}

func (handler *Handler) checkPassword(ctx *gin.Context) {

}

func (handler *Handler) changePassword(ctx *gin.Context) {

}
