package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourorg/social-app/pkg/response"
)

type LikeHandler struct{}

func NewLikeHandler() *LikeHandler {
	return &LikeHandler{}
}

type ToggleLikeRequest struct {
	Action string `json:"action" binding:"required"` // "like" or "unlike"
}

func (h *LikeHandler) Toggle(c *gin.Context) {
	var req ToggleLikeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	// TODO: Implement toggle like
	response.Success(c, gin.H{
		"is_liked":   req.Action == "like",
		"like_count": 0,
	})
}
