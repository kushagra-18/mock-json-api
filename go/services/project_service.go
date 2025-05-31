package services

import (
	"fmt"

	"gorm.io/gorm"
	"mockapi/models" // Assuming module name is mockapi
	"mockapi/utils"  // For GenerateRandomString or similar if needed for ChannelID
)

// ProjectService handles business logic related to projects.
type ProjectService struct {
	DB *gorm.DB
}

// NewProjectService creates a new ProjectService.
func NewProjectService(db *gorm.DB) *ProjectService {
	return &ProjectService{DB: db}
}

// GetProjectByID retrieves a project by its ID.
func (s *ProjectService) GetProjectByID(id uint) (*models.Project, error) {
	var project models.Project
	if err := s.DB.First(&project, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("project with ID %d not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to retrieve project with ID %d: %w", id, err)
	}
	return &project, nil
}

// GetProjectBySlug retrieves a project by its slug.
func (s *ProjectService) GetProjectBySlug(slug string) (*models.Project, error) {
	var project models.Project
	if err := s.DB.Where("slug = ?", slug).First(&project).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("project with slug '%s' not found: %w", slug, err)
		}
		return nil, fmt.Errorf("failed to retrieve project with slug '%s': %w", slug, err)
	}
	return &project, nil
}

// GetProjectByTeamSlugAndProjectSlug retrieves a project by its team's slug and its own slug.
func (s *ProjectService) GetProjectByTeamSlugAndProjectSlug(teamSlug, projectSlug string) (*models.Project, error) {
	var project models.Project
	err := s.DB.Joins("JOIN teams ON teams.id = projects.team_id").
		Where("teams.slug = ? AND projects.slug = ?", teamSlug, projectSlug).
		First(&project).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("project with slug '%s' under team '%s' not found: %w", projectSlug, teamSlug, err)
		}
		return nil, fmt.Errorf("failed to retrieve project '%s' under team '%s': %w", projectSlug, teamSlug, err)
	}
	return &project, nil
}

// CreateProject creates a new project.
// It can set default values, e.g., for ChannelID.
func (s *ProjectService) CreateProject(project *models.Project) error {
	if project.ChannelID == "" {
		// Example: Generate a unique channel ID if not provided
		// This might involve more sophisticated logic in a real application
		randomStr, err := utils.GenerateRandomString(10)
		if err != nil {
			return fmt.Errorf("failed to generate channel ID: %w", err)
		}
		project.ChannelID = "channel_" + randomStr
	}

	if err := s.DB.Create(project).Error; err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}

// UpdateProject updates an existing project.
func (s *ProjectService) UpdateProject(project *models.Project) error {
	if err := s.DB.Save(project).Error; err != nil {
		return fmt.Errorf("failed to update project with ID %d: %w", project.ID, err)
	}
	return nil
}

// UpdateForwardProxyActiveStatus updates the IsForwardProxyActive status of a project.
func (s *ProjectService) UpdateForwardProxyActiveStatus(projectID uint, status bool) error {
	result := s.DB.Model(&models.Project{}).Where("id = ?", projectID).Update("is_forward_proxy_active", status)
	if result.Error != nil {
		return fmt.Errorf("failed to update forward proxy status for project ID %d: %w", projectID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("project with ID %d not found for updating forward proxy status: %w", projectID, gorm.ErrRecordNotFound)
	}
	return nil
}
