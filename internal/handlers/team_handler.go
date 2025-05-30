package handlers

import (
	"errors"
	"go-gin-gorm-api/internal/models"
	"go-gin-gorm-api/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// TeamHandler handles HTTP requests for teams.
type TeamHandler struct {
	teamService services.TeamService
}

// NewTeamHandler creates a new TeamHandler.
func NewTeamHandler(teamService services.TeamService) *TeamHandler {
	return &TeamHandler{teamService: teamService}
}

// CreateTeam handles POST requests to create a new team.
func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var team models.Team
	if err := c.ShouldBindJSON(&team); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Slug generation is handled in the service if team.Slug is empty
	if err := h.teamService.CreateTeam(&team); err != nil {
		// Could check for specific error types, e.g., unique constraint violation
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create team: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, team)
}

// GetAllTeams handles GET requests to retrieve all teams.
func (h *TeamHandler) GetAllTeams(c *gin.Context) {
	teams, err := h.teamService.GetAllTeams()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve teams: " + err.Error()})
		return
	}
	if teams == nil { // Ensure we return an empty list instead of null if no teams
		c.JSON(http.StatusOK, []models.Team{})
		return
	}
	c.JSON(http.StatusOK, teams)
}

// GetTeamByID handles GET requests to retrieve a team by its ID.
func (h *TeamHandler) GetTeamByID(c *gin.Context) {
	idStr := c.Param("teamId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID format"})
		return
	}

	team, err := h.teamService.GetTeamByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // Assuming service might return this
			c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve team: " + err.Error()})
		return
	}
	if team == nil { // Service returns nil if not found (and no other error)
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found"})
		return
	}
	c.JSON(http.StatusOK, team)
}

// GetTeamBySlug handles GET requests to retrieve a team by its slug.
func (h *TeamHandler) GetTeamBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug parameter is required"})
		return
	}

	team, err := h.teamService.GetTeamBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Team not found with slug: " + slug})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve team by slug: " + err.Error()})
		return
	}
	if team == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found with slug: " + slug})
		return
	}
	c.JSON(http.StatusOK, team)
}

// UpdateTeam handles PUT requests to update an existing team.
func (h *TeamHandler) UpdateTeam(c *gin.Context) {
	idStr := c.Param("teamId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID format"})
		return
	}

	var teamUpdates models.Team
	if err := c.ShouldBindJSON(&teamUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	teamUpdates.ID = uint(id) // Set the ID from the path parameter

	// Retrieve existing team to ensure it exists before update
	existingTeam, err := h.teamService.GetTeamByID(uint(id))
	if err != nil || existingTeam == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Team not found to update"})
		return
	}

	// Slug regeneration is handled in the service if teamUpdates.Slug is empty and name changes
	if err := h.teamService.UpdateTeam(&teamUpdates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update team: " + err.Error()})
		return
	}

	// Fetch the updated team to return the full, potentially modified (e.g. slug) object
	updatedTeam, err := h.teamService.GetTeamByID(uint(id))
	if err != nil || updatedTeam == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated team details"})
        return
	}

	c.JSON(http.StatusOK, updatedTeam)
}

// DeleteTeam handles DELETE requests to delete a team by its ID.
func (h *TeamHandler) DeleteTeam(c *gin.Context) {
	idStr := c.Param("teamId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID format"})
		return
	}

	// Optional: Check if team exists before attempting delete
	// existingTeam, err := h.teamService.GetTeamByID(uint(id))
	// if err != nil || existingTeam == nil {
	// 	c.JSON(http.StatusNotFound, gin.H{"error": "Team not found to delete"})
	// 	return
	// }

	if err := h.teamService.DeleteTeam(uint(id)); err != nil {
		// Check if it's a "not found" type of error if GORM returns specific error on delete of non-existent
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete team: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Team deleted successfully"})
}

// RegisterTeamRoutes registers the routes for team operations.
func RegisterTeamRoutes(routerGroup *gin.RouterGroup, teamHandler *TeamHandler) {
	teams := routerGroup.Group("/teams") // Create a subgroup for /teams
	{
		teams.POST("", teamHandler.CreateTeam)
		teams.GET("", teamHandler.GetAllTeams)
		teams.GET("/:teamId", teamHandler.GetTeamByID)
		teams.GET("/slug/:slug", teamHandler.GetTeamBySlug)
		teams.PUT("/:teamId", teamHandler.UpdateTeam)
		teams.DELETE("/:teamId", teamHandler.DeleteTeam)
	}
}
