package routes

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"google.golang.org/genai" // For genai.Client and genai.ClientConfig
	"gorm.io/gorm"

	"mockapi/config"
	"mockapi/controllers"
	"mockapi/middleware"
	"mockapi/services"
)

// SetupRoutes initializes all services, controllers, and sets up the Gin router.
func SetupRoutes(cfg config.Config, db *gorm.DB, redisClient *redis.Client) *gin.Engine {
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Adjust for production
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Initialize Services
	teamService := services.NewTeamService()
	projectService := services.NewProjectService(db)
	randomWordsService := services.NewRandomWordsService()
	requestLogService := services.NewRequestLogService(db)

	if redisClient == nil {
		log.Fatal("Redis client is not initialized for routes setup.")
	}
	redisService := services.NewRedisService(redisClient)

	urlService := services.NewURLService(db, redisService)
	mockContentService := services.NewMockContentService(db)
	proxyService := services.NewProxyService(db)
	fakerService := services.NewFakerService(cfg) // Initialize FakerService

	// AI Prompting Service and Controller
	geminiAPIKey := cfg.GeminiAPIKey
	var aiPromptService services.AIPromptServiceInterface
	var genaiClient *genai.Client

	if geminiAPIKey != "" {
		var errClient error
		genaiClient, errClient = genai.NewClient(context.Background(), &genai.ClientConfig{
			APIKey:  geminiAPIKey,
			Backend: genai.BackendGeminiAPI,
		})
		if errClient != nil {
			log.Printf("Warning: Failed to create genai.Client for AIPromptService: %v. AI features may be disabled.", errClient)
		} else {
			var errService error
			aiPromptService, errService = services.NewAIPromptService(cfg, genaiClient.Models)
			if errService != nil {
				log.Fatalf("Failed to initialize AIPromptService: %v", errService)
			}
			// Note: genaiClient.Close() is handled in main.go
		}
	} else {
		log.Println("GEMINI_API_KEY is not configured. AIPromptService is not initialized.")
	}

	aiPromptController := controllers.NewAIPromptController(aiPromptService)


	// --- Define Routes ---
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "timestamp": time.Now()})
	})

	apiV1 := router.Group("/api/v1")
	{
		homeController := controllers.NewHomeController()
		apiV1.GET("/home", homeController.Home)

		if aiPromptService != nil { // Only register AI route if service is available
			apiV1.POST("/ai/prompt", aiPromptController.HandleAIPrompt)
		} else {
			apiV1.POST("/ai/prompt", func(c *gin.Context) {
				c.JSON(http.StatusServiceUnavailable, gin.H{
					"status": "error",
					"message": "AI service is not configured or available.",
				})
			})
		}
	}

	teamController := controllers.NewTeamController(teamService)
	router.GET("/team", teamController.GetTeamInfo)

	projectController := controllers.NewProjectController(projectService, randomWordsService, requestLogService, teamService)
	projectRoutes := router.Group("/project")
	{
		projectRoutes.POST("/free", projectController.CreateFreeProject)
		projectRoutes.POST("/free/fast-forward", projectController.CreateFreeFastForwardProject)
		projectRoutes.GET("/:projectSlug", projectController.GetProjectBySlug)
	}

	urlController := controllers.NewURLController(urlService)
	urlRoutes := router.Group("/url")
	{
		urlRoutes.PATCH("/:urlId", urlController.UpdateURLInfo)
		urlRoutes.GET("/:urlId", urlController.GetURLDetails)
	}

	proxyController := controllers.NewProxyController(proxyService, projectService)
	proxyRoutes := router.Group("/proxy")
	{
		proxyRoutes.POST("/forward", proxyController.SaveForwardProxy)
		proxyRoutes.PATCH("/forward/active/:projectId", proxyController.UpdateForwardProxyActiveStatus)
	}

	mockContentController := controllers.NewMockContentController(projectService, mockContentService, urlService, requestLogService, redisService, proxyService, fakerService, cfg)
	managementMockRoutes := router.Group("/mock")
	{
		managementMockRoutes.POST("/:projectSlug", mockContentController.SaveMockContent)
		managementMockRoutes.PATCH("/:projectSlug/:urlId", mockContentController.UpdateMockContent)
	}

	authMiddleware := middleware.JWTAuthMiddleware(cfg.JWTSecretKey)
	router.Any("/mock/:teamSlug/:projectSlug/*wildcardPath", authMiddleware, mockContentController.GetMockedJSON)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	return router
}
