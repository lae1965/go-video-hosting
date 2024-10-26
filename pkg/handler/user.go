package handler

import (
	"fmt"
	"go-video-hosting/pkg/model"
	"net/http"
	"strconv"
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
	refreshTokenIdString, err := ctx.Cookie("refreshTokenId")
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("cookie file is missing or defective: %s", err.Error()))
		return
	}

	refreshTokenId, err := strconv.Atoi(refreshTokenIdString)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("cookie file is defective: %s", err.Error()))
		return
	}

	if err := handler.services.Logout(refreshTokenId); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.SetCookie("refreshToken", "", -1, "/", "", false, true)
	ctx.SetCookie("refreshTokenId", "", -1, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{"message": "Logout is success"})
}

func (handler *Handler) refresh(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, fmt.Sprintf("can't find refreshToken in cookie: %s", err.Error()))
		return
	}
	userResponse, err := handler.services.Refresh(refreshToken)
	if err != nil {
		newErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	ctx.SetCookie("refreshToken", userResponse.RefreshToken, int(time.Hour*24*60), "/", "", false, true)
	ctx.SetCookie("refreshTokenId", fmt.Sprint(userResponse.RefreshTokenId), int(time.Hour*24*60), "/", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{
		"id":          userResponse.Id,
		"nickName":    userResponse.NickName,
		"role":        userResponse.Role,
		"isBanned":    userResponse.IsBanned,
		"accessToken": userResponse.AccessToken,
	})
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
	idStr := ctx.Param("id")
	if err := handler.validators.Validate.Var(idStr, "required,numeric,min=1"); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.validators.Validate.Var(file, "avatar"); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.services.SaveAvatar(int(id), file.Filename); err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Avatar was saved successfully"})
}

func (handler *Handler) getAvatar(ctx *gin.Context) {
	idStr := ctx.Param("id")

	if err := handler.validators.Validate.Var(idStr, "required,numeric,min=1"); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	isHeadersSet := false
	err = handler.services.GetAvatar(int(id), func(fileSize int64, mimeType string, chunk []byte) error {
		if !isHeadersSet {
			ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
			ctx.Writer.Header().Set("Content-Type", mimeType)
			ctx.Writer.WriteHeader(http.StatusOK)
			isHeadersSet = true
		}
		_, err := ctx.Writer.Write(chunk)
		return err
	})

	if err != nil {
		newErrorResponse(ctx, http.StatusInternalServerError, err.Error())
	}
}

func (handler *Handler) deleteAvatar(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if err := handler.validators.Validate.Var(idStr, "required,numeric,min=1"); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 0)
	if err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.services.DeleteAvatar(int(id)); err != nil {
		newErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "Avatar was saved successfully"})
}

func (handler *Handler) checkPassword(ctx *gin.Context) {

}

func (handler *Handler) changePassword(ctx *gin.Context) {

}
