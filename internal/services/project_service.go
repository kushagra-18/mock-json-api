package services

import (
	"go-gin-gorm-api/internal/models"
	"go-gin-gorm-api/internal/repositories"
)

// ProjectService defines the interface for project-related business logic.
type ProjectService interface {
	GetProjectByID(id uint) (*models.Project, error)
}

// projectService implements ProjectService.
type projectService struct {
	projectRepo repositories.ProjectRepository
}

// NewProjectService creates a new instance of ProjectService.
func NewProjectService(repo repositories.ProjectRepository) ProjectService {
	return &projectService{projectRepo: repo}
}

// GetProjectByID retrieves a project by its ID.
func (s *projectService) GetProjectByID(id uint) (*models.Project, error) {
	return s.projectRepo.GetProjectByID(id)
}
