package handler

import (
	"fmt"
	"go-video-hosting/internal/errors"
	"go-video-hosting/internal/model"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (handler *Handler) registration(ctx *gin.Context) {
	var input model.Users

	if err := ctx.BindJSON(&input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.validators.Validate.Struct(input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userResponse, err := handler.services.CreateUser(input)
	if err != nil {
		var code int
		switch err.Type {
		case errors.NotUnique:
			code = http.StatusConflict
		default:
			code = http.StatusInternalServerError
		}
		ErrorResponse(ctx, code, err.Message)
		return
	}

	ctx.SetCookie("refreshToken", userResponse.RefreshToken, int(time.Hour*24*60), "/", "", false, true)
	ctx.SetCookie("refreshTokenId", fmt.Sprint(userResponse.RefreshTokenId), int(time.Hour*24*60), "/", "", false, true)

	ctx.JSON(http.StatusCreated, gin.H{"id": userResponse.UserId, "accessToken": userResponse.AccessToken})
}

func (handler *Handler) login(ctx *gin.Context) {
	var input model.Users

	if err := ctx.BindJSON(&input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.validators.Validate.Struct(input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	userResponse, err := handler.services.Login(input)
	if err != nil {
		var code int
		switch err.Type {
		case errors.NotFound:
			code = http.StatusBadRequest
		case errors.Unauthorization:
			code = http.StatusUnauthorized
		default:
			code = http.StatusInternalServerError
		}
		ErrorResponse(ctx, code, err.Message)
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
		ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("cookie file is missing or defective: %s", err.Error()))
		return
	}

	refreshTokenId, err := strconv.Atoi(refreshTokenIdString)
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("cookie file is defective: %s", err.Error()))
		return
	}

	if err := handler.services.Logout(refreshTokenId); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.SetCookie("refreshToken", "", -1, "/", "", false, true)
	ctx.SetCookie("refreshTokenId", "", -1, "/", "", false, true)

	ctx.JSON(http.StatusOK, gin.H{"message": "Logout is success"})
}

func (handler *Handler) refresh(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refreshToken")
	if err != nil {
		ErrorResponse(ctx, http.StatusUnauthorized, fmt.Sprintf("can't find refreshToken in cookie: %s", err.Error()))
		return
	}
	userResponse, appErr := handler.services.Refresh(refreshToken)
	if appErr != nil {
		var code int
		switch appErr.Type {
		case errors.Unauthorization:
			code = http.StatusUnauthorized
		default:
			code = http.StatusInternalServerError
		}

		ErrorResponse(ctx, code, appErr.Message)
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
	id, err := handler.GetIdFromQuery("id", 1, func(key string) string {
		return ctx.Param(key)
	})

	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var jsonObject, user map[string]interface{}
	if err := ctx.BindJSON(&jsonObject); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if updatingObject, exists := jsonObject["updatingObject"]; exists {
		var ok bool
		if user, ok = updatingObject.(map[string]interface{}); !ok {
			ErrorResponse(ctx, http.StatusBadRequest, "invalid updatingObject")
			return
		}
	} else {
		ErrorResponse(ctx, http.StatusBadRequest, `body must contain an object with the key "updatingObject". Data must be inside this object`)
		return
	}

	if email, exist := user["email"].(string); exist {
		if err := handler.validators.Validate.Var(email, "email"); err != nil {
			ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("invalid email format: %s", err.Error()))
			return
		}
	}

	if birthDate, exist := user["birthDate"].(time.Time); exist {
		if _, err := time.Parse(time.RFC3339, birthDate.Format(time.RFC3339)); err != nil {
			ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("invalid birthDate format: %s", err.Error()))
			return
		}
	}

	if appErr := handler.services.UpdateUser(id, user); appErr != nil {
		var code int
		switch appErr.Type {
		case errors.NotUnique:
			code = http.StatusConflict
		case errors.NotFound:
			code = http.StatusNotFound
		default:
			code = http.StatusInternalServerError
		}

		ErrorResponse(ctx, code, appErr.Message)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "update success"})
}

func (handler *Handler) deleteUser(ctx *gin.Context) {
	id, err := handler.GetIdFromQuery("id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.services.DeleteUser(id); err != nil {
		var code int
		switch err.Type {
		case errors.NotFound:
			code = http.StatusNotFound
		default:
			code = http.StatusInternalServerError
		}

		ErrorResponse(ctx, code, err.Message)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (handler *Handler) activate(ctx *gin.Context) {
	activateLink := ctx.Param("link")
	if err := handler.validators.Validate.Var(activateLink, "required,url"); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.services.Activate(activateLink); err != nil {
		var code int
		switch err.Type {
		case errors.NotFound:
			code = http.StatusNotFound
		default:
			code = http.StatusInternalServerError
		}

		ErrorResponse(ctx, code, err.Message)
		return
	}

	ctx.Redirect(http.StatusOK, fmt.Sprintf("%s:%s/emailConfirm", viper.GetString("client.host"), viper.GetString("client.port")))
}

func (handler *Handler) findMin(ctx *gin.Context) {
	id, err := handler.GetIdFromQuery("id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	nickName, appErr := handler.services.FindNickNameById(id)
	if appErr != nil {
		var code int
		switch appErr.Type {
		case errors.NotFound:
			code = http.StatusNotFound
		default:
			code = http.StatusInternalServerError
		}

		ErrorResponse(ctx, code, appErr.Message)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"nickName": nickName})

}

func (handler *Handler) find(ctx *gin.Context) {
	id, err := handler.GetIdFromQuery("id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	user, appErr := handler.services.FindById(id)
	if appErr != nil {
		var code int
		switch appErr.Type {
		case errors.NotFound:
			code = http.StatusNotFound
		default:
			code = http.StatusInternalServerError
		}

		ErrorResponse(ctx, code, appErr.Message)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (handler *Handler) findAll(ctx *gin.Context) {
	users, err := handler.services.FindAll()
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (handler *Handler) saveAvatar(ctx *gin.Context) {
	id, err := handler.GetIdFromQuery("id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.validators.Validate.Var(file, "avatar"); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.services.SaveAvatar(id, file.Filename); err != nil {
		var code int
		switch err.Type {
		case errors.NotFound:
			code = http.StatusNotFound
		default:
			code = http.StatusInternalServerError
		}

		ErrorResponse(ctx, code, err.Message)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Avatar was saved successfully"})
}

func (handler *Handler) getAvatar(ctx *gin.Context) {
	id, err := handler.GetIdFromQuery("id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	isHeadersSet := false
	appErr := handler.services.GetAvatar(id, func(fileSize int64, mimeType string, chunk []byte) error {
		if !isHeadersSet {
			ctx.Writer.Header().Set("Content-Length", fmt.Sprintf("%d", fileSize))
			ctx.Writer.Header().Set("Content-Type", mimeType)
			ctx.Writer.WriteHeader(http.StatusOK)
			isHeadersSet = true
		}
		_, err := ctx.Writer.Write(chunk)
		return err
	})

	if appErr != nil {
		var code int
		switch appErr.Type {
		case errors.NotFound:
			code = http.StatusNotFound
		case errors.EmptyField:
			code = http.StatusNoContent
		default:
			code = http.StatusInternalServerError
		}

		ErrorResponse(ctx, code, appErr.Message)
	}
}

func (handler *Handler) deleteAvatar(ctx *gin.Context) {
	id, err := handler.GetIdFromQuery("id", 1, func(key string) string {
		return ctx.Param(key)
	})
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.services.DeleteAvatar(id); err != nil {
		var code int
		switch err.Type {
		case errors.NotFound:
			code = http.StatusNotFound
		case errors.EmptyField:
			code = http.StatusNotFound
		default:
			code = http.StatusInternalServerError
		}

		ErrorResponse(ctx, code, err.Message)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (handler *Handler) checkUnique(ctx *gin.Context) {
	nickName := ctx.Query("nickName")
	email := ctx.Query("email")

	if nickName == "" && email == "" {
		ErrorResponse(ctx, http.StatusBadRequest, "no objects for check")
		return
	}

	isUnique, message, err := handler.services.CheckIsNickNameEmailUnique(nickName, email)
	if err != nil {
		ErrorResponse(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	if !isUnique {
		ErrorResponse(ctx, http.StatusConflict, message)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (handler *Handler) changePassword(ctx *gin.Context) {
	var input model.ChangePasswordRequest

	if err := ctx.BindJSON(&input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.validators.Validate.Struct(input); err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	refreshTokenIdString, err := ctx.Cookie("refreshTokenId")
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("cookie file is missing or defective: %s", err.Error()))
		return
	}

	refreshTokenId, err := strconv.Atoi(refreshTokenIdString)
	if err != nil {
		ErrorResponse(ctx, http.StatusBadRequest, fmt.Sprintf("cookie file is defective: %s", err.Error()))
		return
	}

	if err := handler.services.ChangePassword(input.Id, refreshTokenId, input.OldPassword, input.NewPassword); err != nil {
		var code int
		switch err.Type {
		case errors.NotFound:
			code = http.StatusNotFound
		case errors.NotEqual:
			code = http.StatusConflict
		default:
			code = http.StatusInternalServerError
		}

		ErrorResponse(ctx, code, err.Message)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "change password success"})
}
