package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHomeHandler_Home(t *testing.T) {
	// Set Gin to TestMode to reduce console output during tests.
	gin.SetMode(gin.TestMode)

	// Create a new Gin router.
	router := gin.New()

	// Instantiate the handler.
	homeHandler := NewHomeHandler()

	// Register the route directly for this unit test.
	// If using RegisterHomeRoutes, the setup would be:
	// apiV1Group := router.Group("/api/v1")
	// RegisterHomeRoutes(apiV1Group, homeHandler)
	// And the request path would be "/api/v1/home"
	// For this test, let's assume the handler is mounted at "/home" for simplicity or directly test its behavior.
	// Let's stick to the specified "/api/v1/home" for consistency with the problem description.
	router.GET("/api/v1/home", homeHandler.Home)

	// Create an HTTP test request.
	req, err := http.NewRequest(http.MethodGet, "/api/v1/home", nil)
	assert.NoError(t, err) // Ensure request creation was successful.

	// Create an HTTP test response recorder.
	rr := httptest.NewRecorder()

	// Serve the HTTP request.
	router.ServeHTTP(rr, req)

	// Assert status code.
	assert.Equal(t, http.StatusOK, rr.Code)

	// Assert response body using JSONEq for robust JSON comparison.
	expectedBody := `{"message":"Hello World"}`
	assert.JSONEq(t, expectedBody, rr.Body.String())
}
