package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"social-app/pkg/response"
)

func HealthCheck(c *gin.Context) {
	response.JSON(c, http.StatusOK, gin.H{
		"code":    0,
		"message": "ok",
		"data": gin.H{
			"status": "healthy",
		},
	})
}
