package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func JSON(c *gin.Context, httpStatus int, data interface{}) {
	c.JSON(httpStatus, data)
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, httpStatus int, code int, message string) {
	c.JSON(httpStatus, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, 30001, message)
}

func Unauthorized(c *gin.Context, message string) {
	Error(c, http.StatusUnauthorized, 10004, message)
}

func Forbidden(c *gin.Context, message string) {
	Error(c, http.StatusForbidden, 20002, message)
}

func NotFound(c *gin.Context, message string) {
	Error(c, http.StatusNotFound, 20001, message)
}

func InternalError(c *gin.Context, message string) {
	Error(c, http.StatusInternalServerError, 30002, message)
}
