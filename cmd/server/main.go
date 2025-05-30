package main

import (
	"fmt"
	"log"

	"go-gin-gorm-api/internal/config"
	_ "go-gin-gorm-api/internal/models" // Ensure models are registered for auto-migration
)

func main() {
	// Replace with your actual DSN
	dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := config.InitDB(dsn)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	if db != nil {
		log.Println("Database connection successful and migrations complete.")
	}
}
