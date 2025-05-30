package main

import (
	"fmt"
	"log"

	"go-gin-gorm-api/internal/config"
	_ "go-gin-gorm-api/internal/models" // Ensure models are registered for auto-migration
	// Gin and other handler/service imports will be added here later
)

func main() {
	// Load application configuration
	appConfig, err := config.LoadConfig(".") // Load from config.yaml in current directory
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	log.Printf("Configuration loaded: %+v", appConfig)


	// Initialize database connection using loaded configuration
	db, err := config.InitDB() // No DSN argument needed now
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if db != nil {
		// This log is now inside InitDB, but can keep one here too
		log.Println("Database connection successfully initialized from main.")
	}

	// Gin router setup and server start will go here later
	// Example:
	// router := gin.Default()
	// router.GET("/ping", func(c *gin.Context) {
	// 	c.JSON(200, gin.H{"message": "pong"})
	// })
	// log.Printf("Starting server on port %s", config.AppConfig.ServerPort)
	// if err := router.Run(":" + config.AppConfig.ServerPort); err != nil {
	// 	log.Fatalf("Failed to start server: %v", err)
	// }
}
