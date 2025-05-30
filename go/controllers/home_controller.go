package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"mockapi/utils" // Module name 'mockapi'
)

// HomeController handles basic home routes.
type HomeController struct {
	// No dependencies for this simple controller
}

// NewHomeController creates a new HomeController.
func NewHomeController() *HomeController {
	return &HomeController{}
}

// Home handles GET /api/v1/home and returns a simple greeting.
func (hc *HomeController) Home(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": "Hello World from MockAPI v1"})
}
