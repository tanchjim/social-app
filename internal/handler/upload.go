package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yourorg/social-app/internal/config"
	"github.com/yourorg/social-app/pkg/response"
)

type UploadHandler struct {
	cfg *config.Config
}

func NewUploadHandler(cfg *config.Config) *UploadHandler {
	return &UploadHandler{cfg: cfg}
}

type SignRequest struct {
	FileName    string `json:"file_name" binding:"required"`
	ContentType string `json:"content_type" binding:"required"`
	FileSize    int64  `json:"file_size" binding:"required"`
}

func (h *UploadHandler) Sign(c *gin.Context) {
	var req SignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	// TODO: Implement COS sign generation
	response.Success(c, gin.H{
		"upload_url":  "",
		"object_key":  "",
		"expired_at":  "",
	})
}
