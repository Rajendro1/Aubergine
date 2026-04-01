package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"aubergine/internal/models"
)

var DB *gorm.DB

func ConnectDB() {
	// For production, these values should come from environment variables.
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Mock DSN for local development
		dsn = "host=localhost user=postgres password=postgres dbname=streamingdb port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("Failed to connect to database: %v. (Ensure Postgres is running if testing locally)", err)
		return
	}

	DB = db
	fmt.Println("Connected to Database successfully.")

	// Auto-migrate tables
	err = DB.AutoMigrate(
		&models.User{},
		&models.Subscription{},
		&models.Video{},
	)
	if err != nil {
		log.Fatalf("AutoMigration failed: %v", err)
	}
}
