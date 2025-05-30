package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mockapi/dtos"
	"mockapi/models"
	"mockapi/services"
	"mockapi/utils"
)

// ProxyController handles forward proxy related API endpoints.
type ProxyController struct {
	proxyService   *services.ProxyService
	projectService *services.ProjectService
}

// NewProxyController creates a new ProxyController.
func NewProxyController(ps *services.ProxyService, prjService *services.ProjectService) *ProxyController {
	return &ProxyController{proxyService: ps, projectService: prjService}
}

// SaveForwardProxy handles POST /proxy/forward
// This will create or update the forward proxy setting for the given ProjectID.
func (pc *ProxyController) SaveForwardProxy(c *gin.Context) {
	var dto dtos.ForwardProxyDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	// Validate that the project exists
	_, err := pc.projectService.GetProjectByID(dto.ProjectID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Project with ID %d not found.", dto.ProjectID))
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error verifying project: "+err.Error())
		}
		return
	}

	proxyModel := &models.ForwardProxy{
		ProjectID: dto.ProjectID,
		Domain:    dto.Domain,
	}

	savedProxy, err := pc.proxyService.CreateForwardProxy(proxyModel, dto.ProjectID)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to save forward proxy settings: "+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, savedProxy) // Or http.StatusCreated if always new
}

// UpdateForwardProxyActiveStatus handles PATCH /proxy/forward/active/:projectId
func (pc *ProxyController) UpdateForwardProxyActiveStatus(c *gin.Context) {
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid Project ID format.")
		return
	}

	var dto dtos.UpdateForwardProxyStatusDTO
	// It's a PATCH, so if body is empty, it might mean no change or use defaults.
	// Binding may error if Content-Type is set but body is empty.
	// For a simple boolean toggle, sometimes it's passed as query param or path param.
	// Here, assuming it comes from JSON body.
	if err := c.ShouldBindJSON(&dto); err != nil {
		// If DTO is empty, isActive will be false. This might be desired default.
		// Depending on strictness, one might allow empty body for PATCH.
		// For now, let's assume an explicit JSON true/false is sent.
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	// Check if project exists
	_, err = pc.projectService.GetProjectByID(uint(projectID))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Project with ID %d not found.", projectID))
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Error verifying project: "+err.Error())
		}
		return
	}

	// Check if forward proxy is configured for the project.
	// It might not be an error if it's not, just that the active status cannot be set.
	// However, the Java code seems to update Project.isForwardProxyActive directly.
	// proxySettings, err := pc.proxyService.GetForwardProxyByProjectID(uint(projectID))
	// if err != nil {
	// 	utils.ErrorResponse(c, http.StatusInternalServerError, "Error fetching proxy settings: "+err.Error())
	// 	return
	// }
	// if proxySettings == nil && dto.IsActive {
	// 	utils.ErrorResponse(c, http.StatusPreconditionFailed, "Forward proxy is not configured for this project. Cannot activate.")
	// 	return
	// }


	if err := pc.projectService.UpdateForwardProxyActiveStatus(uint(projectID), dto.IsActive); err != nil {
		if err == gorm.ErrRecordNotFound { // Service returns this if project not found by Update
			utils.ErrorResponse(c, http.StatusNotFound, fmt.Sprintf("Project with ID %d not found for status update.", projectID))
		} else {
			utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update forward proxy active status: "+err.Error())
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, gin.H{"message": fmt.Sprintf("Forward proxy active status for project %d updated to %v.", projectID, dto.IsActive)})
}
