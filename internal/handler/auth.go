package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/social-app/internal/config"
	"github.com/yourorg/social-app/pkg/response"
)

type AuthHandler struct {
	cfg *config.Config
}

func NewAuthHandler(cfg *config.Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Nickname string `json:"nickname"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// TODO: Implement registration logic
	response.Success(c, gin.H{
		"user_id":   0,
		"username":  req.Username,
		"nickname":  req.Nickname,
	})
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// TODO: Implement login logic
	response.Success(c, gin.H{
		"user_id":       0,
		"username":      req.Username,
		"nickname":      "",
		"avatar":        "",
		"access_token":  "",
		"refresh_token": "",
	})
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// TODO: Implement refresh logic
	response.Success(c, gin.H{
		"access_token":  "",
		"refresh_token": "",
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// TODO: Implement logout logic (revoke all refresh tokens)
	response.Success(c, nil)
}
