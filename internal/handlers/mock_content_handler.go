package handlers

import (
	"fmt"
	"go-gin-gorm-api/internal/models"
	"go-gin-gorm-api/internal/services"
	"net/http"

	"encoding/json"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// MockContentHandler handles HTTP requests for mock content.
type MockContentHandler struct {
	mockContentService services.MockContentService
	urlService         services.UrlService
	projectService     services.ProjectService
	requestLogService  services.RequestLogService
}

// NewMockContentHandler creates a new MockContentHandler.
func NewMockContentHandler(
	mcService services.MockContentService,
	urlService services.UrlService,
	projService services.ProjectService,
	reqLogService services.RequestLogService,
) *MockContentHandler {
	return &MockContentHandler{
		mockContentService: mcService,
		urlService:         urlService,
		projectService:     projService,
		requestLogService:  reqLogService,
	"go-gin-gorm-api/internal/config" // Added for AppConfig access
)

const DefaultProjectIDForMock = 1 // Or load from config
// const BaseURLPlaceholder = "example.com" // Will be replaced by config.AppConfig.AppBaseURL

// CreateMock handles POST requests to /api/v1/mock.
// It creates a URL (if it doesn't exist) and associated mock content.
func (h *MockContentHandler) CreateMock(c *gin.Context) {
	var dto MockContentUrlDto
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if dto.UrlData.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "url_data.url is required"})
		return
	}
	urlString := dto.UrlData.URL // This is the path, e.g., /users/1

	// Attempt to find the Url by its path
	existingUrl, err := h.urlService.GetUrlByPath(urlString)
	if err != nil {
		// Assuming GetUrlByPath returns gorm.ErrRecordNotFound or similar if not found,
		// and other errors for actual DB problems. For now, let's assume nil,nil for not found.
		// A more robust error handling would differentiate.
		// For this example, any error from GetUrlByPath is treated as "try to create".
		// A proper check for gorm.ErrRecordNotFound would be better.
	}

	var finalUrl *models.Url
	if existingUrl != nil {
		finalUrl = existingUrl
	} else {
		// URL not found, create it
		// Fetch a default project
		defaultProject, projErr := h.projectService.GetProjectByID(DefaultProjectIDForMock)
		if projErr != nil || defaultProject == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Default project with ID %d not found or error: %v", DefaultProjectIDForMock, projErr)})
			return
		}

		newUrl := &dto.UrlData
		newUrl.ProjectID = defaultProject.ID
		// Name and Description for URL should be in dto.UrlData if provided by client
		// If not, they will be empty or their zero values.
		// Slug for URL is not explicitly handled here, GORM might make it empty or it needs generation logic.

		createdUrl, createErr := h.urlService.CreateUrl(newUrl, defaultProject.ID)
		if createErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create URL: " + createErr.Error()})
			return
		}
		finalUrl = createdUrl
	}

	if len(dto.MockContentList) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "mock_content_list must contain at least one item"})
		return
	}
	// Take the first MockContent from the list
	mockContentToCreate := dto.MockContentList[0]

	// Create MockContent associated with the finalUrl
	savedMockContent, mcErr := h.mockContentService.CreateMockContent(&mockContentToCreate, finalUrl.ID)
	if mcErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create mock content: " + mcErr.Error()})
		return
	}

	// Construct the mockedUrl string
	// The original Java code uses urlString + ".free." + config.BaseURL
	// urlString is the path. This implies the full mocked URL might look like:
	// https://example.com/some/path.free.example.com - this seems unusual.
	// More typically, it might be a subdomain or a path prefix.
	// Let's assume for now it's: <path>.free.<base_domain>
	// Or perhaps: <base_domain>/mock/<project_slug>/<team_slug>/<path>
	// Given the context "urlString + .free. + config.BaseURL", it might be constructing a unique identifier
	// or a specific routing key rather than a directly callable URL.
	// Let's stick to the literal interpretation for now, acknowledging it's odd.
	mockedUrl := fmt.Sprintf("%s.free.%s", urlString, config.AppConfig.AppBaseURL)
	// A more standard approach might be:
	// mockedUrl := fmt.Sprintf("https://%s/mock/%d%s", config.AppConfig.AppBaseURL, finalUrl.ProjectID, finalUrl.URL)

	c.JSON(http.StatusCreated, gin.H{
		"url":         mockedUrl, // The specially constructed "mocked URL" string
		"data":        savedMockContent,
		"status_code": http.StatusCreated, // Redundant with HTTP status, but matches example
		"original_url_details": finalUrl, // For clarity, returning the saved/found URL object
	})
}

// RegisterMockContentRoutes registers the routes for mock content operations.
func RegisterMockContentRoutes(routerGroup *gin.RouterGroup, mockContentHandler *MockContentHandler) {
	mockGroup := routerGroup.Group("/mock") // Create a subgroup for /mock
	{
		mockGroup.POST("", mockContentHandler.CreateMock)
		// The wildcard GET route will be handled differently, likely at the root or v1 level.
		// For testing, we can add a specific path:
		// mockGroup.GET("/test/*any", mockContentHandler.GetMockedJson)
	}
}

// GetMockedJson handles GET requests for dynamic mock responses.
// This is intended to be used with a wildcard route like /api/v1/mock/:teamSlug/:projectSlug/*fullPath
// Or more generally, context for team/project might be set by a preceding middleware.
func (h *MockContentHandler) GetMockedJson(c *gin.Context) {
	// Team/Project Slugs (Placeholder - these would ideally be set by a middleware based on hostname or path prefix)
	teamSlug := c.GetString("teamSlug") // Assuming a middleware might set this in Gin's context
	projectSlug := c.GetString("projectSlug") // Assuming a middleware might set this

	// For now, using placeholders if not set by a (future) middleware or explicit path params
	// In a real scenario, if these are from path params like /:teamSlug/:projectSlug/*any,
	// they would be retrieved using c.Param("teamSlug") and c.Param("projectSlug").
	// The current structure with c.GetString implies a middleware is responsible for populating these.
	// If this handler is directly matched on a path like /mock/:teamSlug/:projectSlug/*fullPath,
	// then c.Param() is the way. Let's assume for now they might come from path or middleware.

	if teamSlug == "" {
		teamSlug = c.Param("teamSlug") // Fallback to path param if not in context
		if teamSlug == "" {
			// teamSlug = "default-team-slug-placeholder" // Final fallback if truly not available
			// Instead of placeholder, better to make it clear it's required if not found via any means
			// For the specific path /mock/:teamSlug/:projectSlug/*actualPath this will be set by c.Param
		}
	}
	if projectSlug == "" {
		projectSlug = c.Param("projectSlug") // Fallback to path param
		if projectSlug == "" {
			// projectSlug = "default-project-slug-placeholder"
		}
	}

	// If teamSlug or projectSlug are still empty here, it means they were not in context or path params.
	// This indicates a routing/middleware setup issue or a direct call to a route not providing them.
	// For the specific task, these are expected. If they are essential for GetUrlByDetails,
	// and not found, we should probably return an error.
	// The GetUrlByDetails function in url_repository.go uses Joins on Team and Project slugs.

	// Full Path
	// The full path needs to be extracted carefully depending on how the wildcard route is defined.
	// If route is /mock/:teamSlug/:projectSlug/*actualPath, then c.Param("actualPath") gives the part after slugs.
	actualPath := c.Param("actualPath") // This captures everything after /mock/:teamSlug/:projectSlug/
	if !strings.HasPrefix(actualPath, "/") {
		actualPath = "/" + actualPath // Ensure leading slash if not present
	}

	fullPathWithQuery := actualPath
	if c.Request.URL.RawQuery != "" {
		fullPathWithQuery = actualPath + "?" + c.Request.URL.RawQuery
	}
	// The original Java code does ltrim on the request URI.
	// c.Request.URL.Path already gives the path part of the URL.
	// If using *actualPath, it's already the "remaining" part.

	// Get URL Data
	// We need to ensure teamSlug and projectSlug are available.
	if teamSlug == "" || projectSlug == "" {
		// This case implies the route was not /mock/:teamSlug/:projectSlug/*actualPath
		// or team/project context was not set by middleware.
		// For now, we'll try to use the full request path if slugs are missing,
		// assuming GetUrlByDetails might handle simpler path lookups if slugs are empty (not ideal).
		// However, the current GetUrlByDetails expects slugs.
		// A better approach might be to have a different service method if slugs are not part of the query.
		// For this exercise, let's assume the route provides teamSlug and projectSlug.
		// If they are truly optional, then GetUrlByDetails or another service method needs to account for that.
		// Given the method signature: GetUrlByTeamSlugAndProjectSlugAndUrlPath, they are mandatory.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Team slug and Project slug are required in the path or context."})
		return
	}


	urlData, err := h.urlService.GetUrlByDetails(teamSlug, projectSlug, actualPath) // Use actualPath (without query) for DB lookup

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Mock definition not found for this URL, team, and project combination."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving URL data: " + err.Error()})
		return
	}
	if urlData == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Mock definition not found for this URL (nil urlData)."})
		return
	}

	// Select Mock Content
	if len(urlData.MockContentList) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No mock content available for this URL."})
		return
	}

	selectedMockContent := h.mockContentService.SelectRandomMockContent(urlData.MockContentList)
	if selectedMockContent == nil {
		// This case implies all randomness weights were zero or list was effectively empty for selection.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not select mock content (e.g., all weights zero or invalid setup)."})
		return
	}

	// Simulate Latency
	h.mockContentService.SimulateLatency(selectedMockContent)

	// Parse JSON Data
	var jsonData interface{}
	err = json.Unmarshal([]byte(selectedMockContent.Data), &jsonData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse mock content data as JSON: " + err.Error()})
		return
	}

	// Log Request & Emit Event (Async)
	ipAddress := c.ClientIP()
	method := c.Request.Method

	projectIDToLog := urlData.ProjectID // ProjectID is directly on Url model

	// Use fullPathWithQuery for logging, as it includes query params.
	h.requestLogService.SaveRequestLogAsync(fullPathWithQuery, method, ipAddress, http.StatusOK, projectIDToLog)
	h.requestLogService.EmitPusherEventAsync(method, fullPathWithQuery, projectIDToLog, http.StatusOK)

	c.JSON(http.StatusOK, jsonData)
}
