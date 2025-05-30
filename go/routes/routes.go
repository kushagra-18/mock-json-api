package routes

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"

	"mockapi/config"
	"mockapi/controllers"
	"mockapi/middleware" // Import the middleware package
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
	// Note: Order matters if services depend on each other, though constructors usually just take instances.
	teamService := services.NewTeamService(/* db */) // DB part commented out in service
	projectService := services.NewProjectService(db)
	randomWordsService := services.NewRandomWordsService()
	requestLogService := services.NewRequestLogService(db /*, pusherService */) // PusherService not implemented

	// Ensure Redis client is valid before creating services that use it
	if redisClient == nil {
		log.Fatal("Redis client is not initialized. Cannot create services dependent on Redis.")
	}
	redisService := services.NewRedisService(redisClient)

	urlService := services.NewURLService(db, redisService)
	mockContentService := services.NewMockContentService(db)
	proxyService := services.NewProxyService(db)

	// Initialize Controllers
	homeController := controllers.NewHomeController()
	teamController := controllers.NewTeamController(teamService)
	projectController := controllers.NewProjectController(projectService, randomWordsService, requestLogService, teamService)
	urlController := controllers.NewURLController(urlService)
	proxyController := controllers.NewProxyController(proxyService, projectService)
	mockContentController := controllers.NewMockContentController(projectService, mockContentService, urlService, requestLogService, redisService, proxyService, cfg)

	// --- Define Routes ---

	// Simple health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "timestamp": time.Now()})
	})

	// API v1 Group
	apiV1 := router.Group("/api/v1")
	{
		// Home routes
		apiV1.GET("/home", homeController.Home)

		// Other v1 routes can go here if any were strictly under /api/v1 in Java
	}

	// Ungrouped or differently grouped routes (matching Java structure)
	// If Java paths were like "/team", "/project/free", they are top-level under the base URL.

	// Team routes
	// Assuming /team was a base path in Java.
	// The subtask description for TeamController.GetTeamInfo was GET /team.
	router.GET("/team", teamController.GetTeamInfo) // Placeholder, might be /team/:teamSlug or similar

	// Project routes
	projectRoutes := router.Group("/project")
	// JWT Auth middleware could be applied here: projectRoutes.Use(middlewares.AuthMiddleware(cfg.JWTSecretKey))
	{
		projectRoutes.POST("/free", projectController.CreateFreeProject)
		projectRoutes.POST("/free/fast-forward", projectController.CreateFreeFastForwardProject)
		projectRoutes.GET("/:projectSlug", projectController.GetProjectBySlug)
		// Example: projectRoutes.GET("/:projectSlug/details", middlewares.AuthMiddleware(cfg.JWTSecretKey), projectController.GetProjectDetails)
	}

	// URL routes (assuming /url base path)
	urlRoutes := router.Group("/url")
	// urlRoutes.Use(middlewares.AuthMiddleware(cfg.JWTSecretKey))
	{
		urlRoutes.PATCH("/:urlId", urlController.UpdateURLInfo) // Should be PATCH for update
		urlRoutes.GET("/:urlId", urlController.GetURLDetails)
	}

	// Proxy routes
	proxyRoutes := router.Group("/proxy")
	// proxyRoutes.Use(middlewares.AuthMiddleware(cfg.JWTSecretKey))
	{
		proxyRoutes.POST("/forward", proxyController.SaveForwardProxy)
		proxyRoutes.PATCH("/forward/active/:projectId", proxyController.UpdateForwardProxyActiveStatus)
	}

	// Mock Content serving and management routes
	// Management (POST, PATCH for definitions) - potentially needs auth
	managementMockRoutes := router.Group("/mock")
	// managementMockRoutes.Use(middlewares.AuthMiddleware(cfg.JWTSecretKey))
	{
		managementMockRoutes.POST("/:projectSlug", mockContentController.SaveMockContent)
		managementMockRoutes.PATCH("/:projectSlug/:urlId", mockContentController.UpdateMockContent)
	}

	// Public serving endpoint for mocks
	// Path: /mock/:teamSlug/:projectSlug/*wildcardPath
	// This endpoint specifically requires JWT authentication as per Java's @RequestParam token.
	authMiddleware := middleware.JWTAuthMiddleware(cfg.JWTSecretKey)

	// Apply auth middleware to the GetMockedJSON route.
	// Using .ANY to match all HTTP methods as the original Java controller method
	// did not specify a particular HTTP method (implying it could handle multiple or defaulted to GET).
	// If specific methods are known, router.GET, router.POST etc. can be used with the middleware.
	router.Any("/mock/:teamSlug/:projectSlug/*wildcardPath", authMiddleware, mockContentController.GetMockedJSON)


	// Fallback for undefined routes (404)
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	return router
}
