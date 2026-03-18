package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/social-app/internal/config"
	"github.com/yourorg/social-app/pkg/jwt"
	"github.com/yourorg/social-app/pkg/response"
)

func Auth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Missing Authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid Authorization header format")
			c.Abort()
			return
		}

		claims, err := jwt.ValidateToken(parts[1], cfg.JWTSecret)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, 10004, "Token 无效")
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Next()
	}
}
