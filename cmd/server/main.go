package main

import (
	"log"

	"social-app/internal/config"
	"social-app/internal/router"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize router
	r := router.Setup(cfg)

	// Start server
	log.Printf("Server starting on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
