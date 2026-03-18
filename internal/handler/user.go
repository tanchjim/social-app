package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourorg/social-app/pkg/response"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement get user by ID
	response.Success(c, gin.H{
		"user_id":       id,
		"username":      "",
		"nickname":      "",
		"avatar":        "",
		"bio":           "",
		"content_count": 0,
		"like_count":    0,
		"created_at":    "",
	})
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Bio      string `json:"bio"`
	Avatar   string `json:"avatar"`
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	// TODO: Implement update user
	response.Success(c, gin.H{
		"user_id":  id,
		"nickname": req.Nickname,
		"bio":      req.Bio,
		"avatar":   req.Avatar,
	})
}

func (h *UserHandler) GetMyContents(c *gin.Context) {
	// TODO: Implement get my contents
	response.Success(c, gin.H{
		"list":      []interface{}{},
		"total":     0,
		"page":      1,
		"page_size": 20,
	})
}
