package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Message   string `json:"message,omitempty"`
	Data      any    `json:"data,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
}

func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, APIResponse{
		Data: data,
	})
}

func SuccessWithMessage(c *gin.Context, data any) {
	c.JSON(http.StatusOK, APIResponse{
		Data: data,
	})
}

func Created(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, data)
}

func Error(c *gin.Context, httpStatus int, message string, errorCode string) {
	c.AbortWithStatusJSON(httpStatus, APIResponse{
		ErrorCode: errorCode,
		Message:   message,
	})
}

func ValidationError(c *gin.Context, errors string) {
	c.AbortWithStatusJSON(http.StatusUnprocessableEntity, APIResponse{
		Message:   "Validation failed",
		ErrorCode: errors,
	})
}

func BadRequest(c *gin.Context, message, errorCode string) {
	Error(c, http.StatusBadRequest, message, errorCode)
}

func Unauthorized(c *gin.Context, message, errorCode string) {
	Error(c, http.StatusUnauthorized, message, errorCode)
}

func Forbidden(c *gin.Context, message, errorCode string) {
	Error(c, http.StatusForbidden, message, errorCode)
}

func NotFound(c *gin.Context, message, errorCode string) {
	Error(c, http.StatusNotFound, message, errorCode)
}

func InternalServerError(c *gin.Context, message, errorCode string) {
	Error(c, http.StatusInternalServerError, message, errorCode)
}
