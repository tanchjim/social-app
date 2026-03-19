package handler

import (
	"github.com/gin-gonic/gin"
	"social-app/pkg/response"
)

type ContentHandler struct{}

func NewContentHandler() *ContentHandler {
	return &ContentHandler{}
}

type CreateContentRequest struct {
	Type        string   `json:"type" binding:"required"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	MediaURL    string   `json:"media_url" binding:"required"`
	CoverURL    string   `json:"cover_url"`
	Tags        []string `json:"tags"`
}

func (h *ContentHandler) Create(c *gin.Context) {
	var req CreateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	// TODO: Implement create content
	response.Success(c, gin.H{
		"content_id": 0,
		"status":     "published",
		"type":       req.Type,
		"title":      req.Title,
		"media_url":  req.MediaURL,
		"cover_url":  req.CoverURL,
		"created_at": "",
	})
}

func (h *ContentHandler) List(c *gin.Context) {
	// TODO: Implement list contents
	response.Success(c, gin.H{
		"list":      []interface{}{},
		"total":     0,
		"page":      1,
		"page_size": 20,
	})
}

func (h *ContentHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement get content by ID
	response.Success(c, gin.H{
		"content_id":    id,
		"type":          "",
		"title":         "",
		"description":   "",
		"status":        "published",
		"cover_url":     "",
		"media_url":     "",
		"author":        gin.H{},
		"like_count":    0,
		"comment_count": 0,
		"is_liked":      false,
		"tags":          []string{},
		"created_at":    "",
	})
}

func (h *ContentHandler) GetReviewResult(c *gin.Context) {
	id := c.Param("id")
	// TODO: Implement get review result
	response.Success(c, gin.H{
		"content_id":   id,
		"status":       "published",
		"reject_reason": "",
		"reviewed_at":  "",
	})
}

func (h *ContentHandler) Delete(c *gin.Context) {
	// TODO: Implement delete content (soft delete)
	response.Success(c, nil)
}
