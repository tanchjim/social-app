package router

import (
	"github.com/gin-gonic/gin"
	"social-app/internal/config"
	"social-app/internal/handler"
	"social-app/internal/middleware"
)

func Setup(cfg *config.Config) *gin.Engine {
	gin.SetMode(cfg.GinMode)
	r := gin.New()

	// Middlewares
	r.Use(middleware.Logger())
	r.Use(middleware.CORS())
	r.Use(gin.Recovery())

	// Handlers
	authHandler := handler.NewAuthHandler(cfg)
	userHandler := handler.NewUserHandler()
	contentHandler := handler.NewContentHandler()
	commentHandler := handler.NewCommentHandler()
	likeHandler := handler.NewLikeHandler()
	uploadHandler := handler.NewUploadHandler(cfg)

	// Health check
	r.GET("/health", handler.HealthCheck)

	// API v1
	v1 := r.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.Refresh)
			auth.POST("/logout", middleware.Auth(cfg), authHandler.Logout)
		}

		// Upload routes
		v1.POST("/uploads/sign", middleware.Auth(cfg), uploadHandler.Sign)

		// User routes
		users := v1.Group("/users")
		{
			users.GET("/:id", userHandler.GetByID)
			users.PUT("/:id", middleware.Auth(cfg), userHandler.Update)
			users.GET("/me/contents", middleware.Auth(cfg), userHandler.GetMyContents)
		}

		// Content routes
		contents := v1.Group("/contents")
		{
			contents.POST("", middleware.Auth(cfg), contentHandler.Create)
			contents.GET("", contentHandler.List)
			contents.GET("/:id", contentHandler.GetByID)
			contents.GET("/:id/review-result", middleware.Auth(cfg), contentHandler.GetReviewResult)
			contents.DELETE("/:id", middleware.Auth(cfg), contentHandler.Delete)

			// Comment routes
			contents.GET("/:id/comments", commentHandler.List)
			contents.POST("/:id/comments", middleware.Auth(cfg), commentHandler.Create)

			// Like routes
			contents.POST("/:id/like", middleware.Auth(cfg), likeHandler.Toggle)
		}
	}

	return r
}
