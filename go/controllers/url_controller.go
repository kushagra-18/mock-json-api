package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mockapi/dtos"
	"mockapi/services"
	"mockapi/utils"
)

// URLController handles URL-related API endpoints.
type URLController struct {
	urlService *services.URLService
}

// NewURLController creates a new URLController.
func NewURLController(us *services.URLService) *URLController {
	return &URLController{urlService: us}
}

// UpdateURLInfo handles PATCH /url/:urlId
func (uc *URLController) UpdateURLInfo(c *gin.Context) {
	urlIDStr := c.Param("urlId")
	urlID, err := strconv.ParseUint(urlIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid URL ID format.")
		return
	}

	var dto dtos.URLDataDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	// Ensure at least one field is being updated, or handle empty DTO.
	// For now, service layer's UpdateURL will fetch the URL and update fields present in DTO.

	updatedURL, err := uc.urlService.UpdateURL(uint(urlID), dto)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, fmt.Sprintf("URL with ID %d not found.", urlID))
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update URL: "+err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, updatedURL)
}

// GetURLDetails handles GET /url/:urlId
func (uc *URLController) GetURLDetails(c *gin.Context) {
	urlIDStr := c.Param("urlId")
	urlID, err := strconv.ParseUint(urlIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid URL ID format.")
		return
	}

	urlDetails, err := uc.urlService.GetURLByID(uint(urlID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, fmt.Sprintf("URL with ID %d not found.", urlID))
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to retrieve URL details: "+err.Error())
		}
		return
	}

	// Response should ideally be a DTO that includes URL details and its MockContents.
	// For now, returning the model directly.
	utils.SuccessResponse(c, http.StatusOK, urlDetails)
}
