package errors

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ErrorRes struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func NewErrorResponse(ctx *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	ctx.AbortWithStatusJSON(statusCode, ErrorRes{Message: message})
}
func (errRes *ErrorRes) ErrorResponse(ctx *gin.Context, err *ErrorRes) {
	logrus.Error(err.Message)
	ctx.AbortWithStatusJSON(err.Code, gin.H{"message": err.Message})
}
