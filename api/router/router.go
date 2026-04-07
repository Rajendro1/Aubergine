package router

import (
	"aubergine/api/handlers"
	"aubergine/api/middleware"
	customlogger "aubergine/logger"
	"log"

	"github.com/gin-gonic/gin"
)

func HandleRoutes() *gin.Engine {
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

	r := gin.New()
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

		// Plans (Public)
		api.GET("/plans", handlers.GetPlans)

		// Content Discovery (Public)
		api.GET("/content", handlers.ListContent)
		api.GET("/content/recommendations", handlers.GetContentRecommendations)

		// Protected Routes
		protected := api.Group("/")
		protected.Use(middleware.AuthRequired())
		{
			// Auth Profile
			protected.GET("/auth/profile", handlers.GetProfile)
			protected.PUT("/auth/profile", handlers.UpdateProfile)

			// Subscription / Billing
			protected.POST("/subscriptions/subscribe", handlers.Subscribe)
			protected.GET("/subscriptions/history", handlers.GetSubscriptionHistory)
			protected.POST("/subscriptions/:id/cancel", handlers.CancelSubscription)

			// Content Streaming (Dynamic Auth)
			protected.GET("/content/:id/stream", handlers.StreamContent)

			// History
			protected.POST("/history/progress", handlers.UpdateProgress)
			protected.GET("/history/continue-watching", handlers.GetContinueWatching)

			// Admin Routes
			admin := protected.Group("/admin")
			admin.Use(middleware.AdminRequired())
			{
				// Admin Plans
				admin.POST("/plans", handlers.CreatePlan)
				admin.PUT("/plans/:id", handlers.UpdatePlan)
				admin.DELETE("/plans/:id", handlers.DeletePlan)

				// Admin Content
				admin.POST("/content", handlers.CreateContent)
				admin.PUT("/content/:id", handlers.UpdateContent)
				admin.DELETE("/content/:id", handlers.DeleteContent)
			}
		}
	}

	return r
}
