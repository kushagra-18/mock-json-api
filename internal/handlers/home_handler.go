package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeHandler handles requests related to the home endpoint.
type HomeHandler struct {
	// No dependencies for this simple handler, but can be added if needed.
}

// NewHomeHandler creates a new instance of HomeHandler.
func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

// Home handles GET requests to /api/v1/home.
// It returns a simple JSON message.
func (h *HomeHandler) Home(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Hello World",
	})
}

// RegisterHomeRoutes registers the home routes with the given router group.
func RegisterHomeRoutes(routerGroup *gin.RouterGroup, homeHandler *HomeHandler) {
	routerGroup.GET("/home", homeHandler.Home)
}
