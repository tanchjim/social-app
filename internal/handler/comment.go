package handler

import (
	"github.com/gin-gonic/gin"
	"social-app/pkg/response"
)

type CommentHandler struct{}

func NewCommentHandler() *CommentHandler {
	return &CommentHandler{}
}

func (h *CommentHandler) List(c *gin.Context) {
	// TODO: Implement list comments
	response.Success(c, gin.H{
		"list":      []interface{}{},
		"total":     0,
		"page":      1,
		"page_size": 20,
	})
}

type CreateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *CommentHandler) Create(c *gin.Context) {
	var req CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	// TODO: Implement create comment
	response.Success(c, gin.H{
		"comment_id": 0,
		"content":    req.Content,
		"author":     gin.H{},
		"created_at": "",
	})
}
