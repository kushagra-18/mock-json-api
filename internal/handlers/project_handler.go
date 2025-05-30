package handlers

import (
	"errors"
	"fmt"
	"go-gin-gorm-api/internal/models"
	"go-gin-gorm-api/internal/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const DefaultTeamSlug = "default-team" // Or some other identifier for the free tier team

// ProjectHandler handles HTTP requests for projects.
type ProjectHandler struct {
	projectService services.ProjectService
	teamService    services.TeamService
}

// NewProjectHandler creates a new ProjectHandler.
func NewProjectHandler(projectService services.ProjectService, teamService services.TeamService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService, teamService: teamService}
}

// CreateProject handles POST requests to create a new project.
func (h *ProjectHandler) CreateProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	if project.TeamID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "TeamID is required"})
		return
	}

	// Validate TeamID
	team, err := h.teamService.GetTeamByID(project.TeamID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || team == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Team with ID %d not found", project.TeamID)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to validate team: " + err.Error()})
		return
	}
	if team == nil { // Double check after error check
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Team with ID %d not found", project.TeamID)})
		return
	}


	if err := h.projectService.CreateProject(&project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create project: " + err.Error()})
		return
	}
	// Preload team for the response
	project.Team = team
	c.JSON(http.StatusCreated, project)
}

// CreateFreeProject handles POST requests for creating a project under a default/free team.
func (h *ProjectHandler) CreateFreeProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Fetch the default team
	defaultTeam, err := h.teamService.GetTeamBySlug(DefaultTeamSlug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || defaultTeam == nil {
			// Potentially create the default team if it doesn't exist, or return error
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Default team '%s' not found. Please ensure it exists.", DefaultTeamSlug)})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve default team: " + err.Error()})
		return
	}
	if defaultTeam == nil { // Double check
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Default team '%s' not found after check. Please ensure it exists.", DefaultTeamSlug)})
		return
	}

	project.TeamID = defaultTeam.ID
	project.Team = defaultTeam // For response

	if err := h.projectService.CreateProject(&project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create free project: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, project)
}

// GetAllProjects handles GET requests to retrieve all projects.
func (h *ProjectHandler) GetAllProjects(c *gin.Context) {
	projects, err := h.projectService.GetAllProjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve projects: " + err.Error()})
		return
	}
	if projects == nil {
		c.JSON(http.StatusOK, []models.Project{})
		return
	}
	c.JSON(http.StatusOK, projects)
}

// GetProjectByID handles GET requests to retrieve a project by its ID.
func (h *ProjectHandler) GetProjectByID(c *gin.Context) {
	idStr := c.Param("projectId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	project, err := h.projectService.GetProjectByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project: " + err.Error()})
		return
	}
	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}
	c.JSON(http.StatusOK, project)
}

// GetProjectBySlug handles GET requests to retrieve a project by its slug.
func (h *ProjectHandler) GetProjectBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Slug parameter is required"})
		return
	}
	project, err := h.projectService.GetProjectBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found with slug: " + slug})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve project by slug: " + err.Error()})
		return
	}
	if project == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found with slug: " + slug})
		return
	}
	c.JSON(http.StatusOK, project)
}

// GetProjectsByTeamID handles GET requests to retrieve projects for a specific team.
func (h *ProjectHandler) GetProjectsByTeamID(c *gin.Context) {
	teamIdStr := c.Param("teamId")
	teamId, err := strconv.ParseUint(teamIdStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid team ID format"})
		return
	}

	// Optional: Validate TeamID exists
	team, err := h.teamService.GetTeamByID(uint(teamId))
	if err != nil || team == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Team with ID %d not found", teamId)})
		return
	}

	projects, err := h.projectService.GetProjectsByTeamID(uint(teamId))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve projects for team: " + err.Error()})
		return
	}
	if projects == nil {
		c.JSON(http.StatusOK, []models.Project{})
		return
	}
	c.JSON(http.StatusOK, projects)
}

// UpdateProject handles PUT requests to update an existing project.
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	idStr := c.Param("projectId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	var projectUpdates models.Project
	if err := c.ShouldBindJSON(&projectUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload: " + err.Error()})
		return
	}

	// Ensure the ID from path is used, not from body (if present)
	projectUpdates.ID = uint(id)

	// Validate TeamID if provided in the update.
	// If TeamID is 0 or not present in payload, it means team is not being changed.
	// If TeamID is provided, validate it.
	if projectUpdates.TeamID != 0 {
		team, err := h.teamService.GetTeamByID(projectUpdates.TeamID)
		if err != nil || team == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Team with ID %d not found for project update", projectUpdates.TeamID)})
			return
		}
	} else {
		// If TeamID is not in payload, retain the existing one.
		// Fetch existing project to get its current TeamID.
		existingProject, err := h.projectService.GetProjectByID(uint(id))
		if err != nil || existingProject == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found to update"})
			return
		}
		projectUpdates.TeamID = existingProject.TeamID // Keep original TeamID
	}


	if err := h.projectService.UpdateProject(&projectUpdates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update project: " + err.Error()})
		return
	}

	updatedProject, err := h.projectService.GetProjectByID(uint(id))
    if err != nil || updatedProject == nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve updated project details"})
        return
    }
	c.JSON(http.StatusOK, updatedProject)
}

// DeleteProject handles DELETE requests to delete a project by its ID.
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	idStr := c.Param("projectId")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID format"})
		return
	}

	if err := h.projectService.DeleteProject(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete project: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}

// RegisterProjectRoutes registers the routes for project operations.
func RegisterProjectRoutes(baseRouterGroup *gin.RouterGroup, teamRouterGroup *gin.RouterGroup, projectHandler *ProjectHandler) {
	projectsGroup := baseRouterGroup.Group("/projects")
	{
		projectsGroup.POST("", projectHandler.CreateProject)
		projectsGroup.POST("/create-free", projectHandler.CreateFreeProject)
		projectsGroup.GET("", projectHandler.GetAllProjects)
		projectsGroup.GET("/:projectId", projectHandler.GetProjectByID)
		projectsGroup.GET("/slug/:slug", projectHandler.GetProjectBySlug)
		projectsGroup.PUT("/:projectId", projectHandler.UpdateProject)
		projectsGroup.DELETE("/:projectId", projectHandler.DeleteProject)
	}

	// Register GET /teams/:teamId/projects
	// teamRouterGroup is expected to be the one set up for /teams (e.g., apiV1.Group("/teams"))
	if teamRouterGroup != nil {
		teamRouterGroup.GET("/:teamId/projects", projectHandler.GetProjectsByTeamID)
	}
}
