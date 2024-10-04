package handler

import "fmt"

import (
	"go-video-hosting/pkg/model"
	"net/http"
	"time"

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

	userResponse, err := handler.services.CreateUser(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.SetCookie("refreshToken", userResponse.RefreshToken, int(time.Hour*24*60), "/", "", false, true)
	ctx.SetCookie("refreshTokenId", fmt.Sprint(userResponse.RefreshTokenId), int(time.Hour*24*60), "/", "", false, true)
	ctx.JSON(http.StatusCreated, gin.H{"id": userResponse.UserId, "accessToken": userResponse.AccessToken})
}

func (handler *Handler) login(ctx *gin.Context) {
	var input model.Users

	if err := ctx.BindJSON(&input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.validators.Validate.Struct(input); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userResponse, err := handler.services.Login(input)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.SetCookie("refreshToken", userResponse.RefreshToken, int(time.Hour*24*60), "/", "", false, true)
	ctx.SetCookie("refreshTokenId", fmt.Sprint(userResponse.RefreshTokenId), int(time.Hour*24*60), "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"id":               userResponse.Id,
		"nickName":         userResponse.NickName,
		"email":            userResponse.Email,
		"firstName":        userResponse.FirstName,
		"lastName":         userResponse.LastName,
		"birthDay":         userResponse.BirthDate,
		"createdTimestamp": userResponse.CreateTimestamp,
		"role":             userResponse.Role,
		"isBanned":         userResponse.IsBanned,
		"channelsCount":    userResponse.ChannelsCount,
		"accessToken":      userResponse.AccessToken,
	})
}

func (handler *Handler) logout(ctx *gin.Context) {

}

func (handler *Handler) refresh(ctx *gin.Context) {

}

func (handler *Handler) editUser(ctx *gin.Context) {

}

func (handler *Handler) deleteUser(ctx *gin.Context) {

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
