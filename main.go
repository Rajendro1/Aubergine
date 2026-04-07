package main

import (
	"fmt"
	"log"

	"aubergine/api/router"
	"aubergine/internal/database"

	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file (if exists)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	database.ConnectDB()
}

func main() {
	r := router.HandleRoutes()
	port := ":8080"
	fmt.Printf("Starting Streaming API Server on port %s\n", port)
	if err := r.Run(port); err != nil {
		log.Fatal("Server failed:", err)
	}
}
