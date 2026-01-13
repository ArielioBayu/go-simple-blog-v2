package main

import (
	"log"

	"github.com/ArielioBayu/go-simple-blog-v2/config"
	"github.com/ArielioBayu/go-simple-blog-v2/middleware"
	"github.com/ArielioBayu/go-simple-blog-v2/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Set Gin mode
	gin.SetMode(cfg.Server.GinMode)

	// Connect to database
	if err := config.ConnectDB(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Setup middleware
	r.Use(middleware.CORSMiddleware())

	// Setup routes
	routes.SetupRoutes(r)

	// Start server
	address := ":" + cfg.Server.Port
	log.Printf("Starting server on %s", address)
	if err := r.Run(address); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
