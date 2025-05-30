package services

import (
	"go-gin-gorm-api/internal/models"
	"go-gin-gorm-api/internal/repositories"
)

// ProjectService defines the interface for project-related business logic.
type ProjectService interface {
	CreateProject(project *models.Project) error
	GetProjectByID(id uint) (*models.Project, error)
	GetProjectBySlug(slug string) (*models.Project, error) // Added
	GetProjectsByTeamID(teamID uint) ([]models.Project, error) // Added
	GetAllProjects() ([]models.Project, error) // Added
	UpdateProject(project *models.Project) error // Added
	DeleteProject(id uint) error // Added
}

// projectService implements ProjectService.
type projectService struct {
	projectRepo repositories.ProjectRepository
}

// NewProjectService creates a new instance of ProjectService.
func NewProjectService(repo repositories.ProjectRepository) ProjectService {
	return &projectService{projectRepo: repo}
}

// CreateProject creates a new project.
func (s *projectService) CreateProject(project *models.Project) error {
	// Add any business logic here, e.g., generating a slug if not provided
	if project.Slug == "" && project.Name != "" {
		project.Slug = generateSlug(project.Name) // Assuming generateSlug is accessible or redefined
	}
	return s.projectRepo.CreateProject(project)
}

// GetProjectByID retrieves a project by its ID.
func (s *projectService) GetProjectByID(id uint) (*models.Project, error) {
	return s.projectRepo.GetProjectByID(id)
}

// GetProjectBySlug retrieves a project by its slug.
func (s *projectService) GetProjectBySlug(slug string) (*models.Project, error) {
	return s.projectRepo.GetProjectBySlug(slug)
}

// GetProjectsByTeamID retrieves all projects for a given team ID.
func (s *projectService) GetProjectsByTeamID(teamID uint) ([]models.Project, error) {
	return s.projectRepo.GetProjectsByTeamID(teamID)
}

// GetAllProjects retrieves all projects.
func (s *projectService) GetAllProjects() ([]models.Project, error) {
	return s.projectRepo.GetAllProjects()
}

// UpdateProject updates an existing project.
func (s *projectService) UpdateProject(project *models.Project) error {
	if project.Slug == "" && project.Name != "" {
		project.Slug = generateSlug(project.Name) // Assuming generateSlug is accessible or redefined
	}
	return s.projectRepo.UpdateProject(project)
}

// DeleteProject deletes a project by its ID.
func (s *projectService) DeleteProject(id uint) error {
	return s.projectRepo.DeleteProject(id)
}

// generateSlug is a helper function (can be moved to a common util package)
// For now, keeping it simple and local if not already exposed by teamService.
// If teamService.generateSlug is public, that can be used too.
func generateSlug(name string) string {
	// Basic slug generation, replace with a more robust library if needed
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}

// Need to import "strings" if not already present
import "strings"
