package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"

	"mockapi/config"   // Adjust if your module path is different
	"mockapi/database" // Adjust if your module path is different
	"mockapi/routes"   // Adjust if your module path is different
)

func main() {
	// Load configuration
	// Assuming .env file is in the same directory as the executable, or path is set in LoadConfig
	// If running `go run main.go` from /app/go, and .env is in /app/go, "." is correct.
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	database.ConnectDB(cfg)
	sqlDB, err := database.DB.DB()
	if err != nil {
		log.Fatalf("Failed to get underlying sql.DB from GORM: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
	}()


	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Ping Redis to check connection
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully.")
	defer func() {
		if err := redisClient.Close(); err != nil {
			log.Printf("Error closing Redis client: %v", err)
		}
	}()

	// Set Gin mode
	// Consider making this configurable, e.g., via cfg.GinMode
	gin.SetMode(gin.DebugMode) // Defaulting to DebugMode

	// Setup router
	// Ensure database.DB is the gorm.DB instance
	router := routes.SetupRoutes(cfg, database.DB, redisClient)

	// Define server address
	serverAddr := ":" + cfg.ServerPort
	if cfg.ServerPort == "" {
		log.Println("Warning: ServerPort not set in config, defaulting to :8080")
		serverAddr = ":8080" // Ensure a default if config is missing it and not handled by LoadConfig
	}


	// Create an http.Server instance for graceful shutdown
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: router,
		// Good practice: Set timeouts to avoid Slowloris attacks.
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start the server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort) // Use cfg.ServerPort for the log message
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Implement graceful shutdown
	quit := make(chan os.Signal, 1)
	// signal.Notify registers the given channel to receive notifications of the specified signals.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received.
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}
