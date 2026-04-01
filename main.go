package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"aubergine/api/handlers"
	"aubergine/api/middleware"
	"aubergine/internal/database"
	customlogger "aubergine/logger"
)

func main() {
	// 1. Initialize custom async logger
	cfg := customlogger.Config{
		WorkerCount: 5,
		Brokers:     []string{}, // Add kafka brokers here
		Topic:       "api-logs",
		LogDir:      "logs",
	}

	appLogger, err := customlogger.New(cfg)
	if err != nil {
		log.Fatalf("Logger init failed: %v", err)
	}
	defer appLogger.Close()

	// 2. Connect to Database (AutoMigrates schema as well)
	database.ConnectDB()

	// 3. Setup Gin Router
	r := gin.New() // New() creates empty router without default middlewares
	r.Use(gin.Recovery())
	r.Use(middleware.GinLogger(appLogger)) // Inject custom logger

	// Setup public routes
	api := r.Group("/api/v1")
	{
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
				"status":  "active",
			})
		})

		// Authentication Endpoints
		authGroup := api.Group("/auth")
		{
			authGroup.POST("/register", handlers.Register)
			authGroup.POST("/login", handlers.Login)
		}

		// Content Discovery (Public)
		api.GET("/videos", handlers.ListVideos)

		// Protected Routes
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			// Subscription / Billing
			protected.POST("/billing/subscribe", handlers.CreateCheckoutSession)

			// Content Streaming (Dynamic Auth)
			protected.GET("/videos/:id/stream", handlers.StreamVideo)
		}
	}

	port := ":8080"
	fmt.Printf("Starting Streaming API Server on port %s\n", port)
	
	if err := r.Run(port); err != nil {
		appLogger.Error("Server failed to start", map[string]string{"error": err.Error()})
		log.Fatal("Server failed:", err)
	}
}
