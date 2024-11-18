package handler

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorization = "Authorization"
	bearer        = "Bearer"
	accessKey     = "ACCESS_KEY"
	userIdCtx     = "useId"
)

func (handler Handler) AuthMiddleware(ctx *gin.Context) {
	authHeader := ctx.GetHeader(authorization)
	if authHeader == "" {
		ErrorResponse(ctx, http.StatusUnauthorized, fmt.Sprintf(`header "%s" is empty`, authorization))
		return
	}

	authParts := strings.Split(authHeader, " ")
	if len(authParts) != 2 || authParts[0] != bearer {
		ErrorResponse(ctx, http.StatusUnauthorized, fmt.Sprintf(`invalid format of header "%s"`, authorization))
		return
	}

	if len(authParts[1]) == 0 {
		ErrorResponse(ctx, http.StatusUnauthorized, "token is empty")
		return
	}

	userId, err := handler.services.ValidateToken(authParts[1], os.Getenv(accessKey))
	if err != nil {
		ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	ctx.Set(userIdCtx, userId)
}
