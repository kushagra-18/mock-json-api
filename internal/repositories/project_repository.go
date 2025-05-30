package repositories

import (
	"errors"
	"go-gin-gorm-api/internal/models"

	"gorm.io/gorm"
)

// ProjectRepository defines the interface for project data operations.
type ProjectRepository interface {
	CreateProject(project *models.Project) error
	GetProjectByID(id uint) (*models.Project, error)
	GetProjectBySlug(slug string) (*models.Project, error)
	GetProjectsByTeamID(teamID uint) ([]models.Project, error)
	UpdateProject(project *models.Project) error
	DeleteProject(id uint) error
	GetAllProjects() ([]models.Project, error)
}

// projectRepository implements ProjectRepository with GORM.
type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new instance of projectRepository.
func NewProjectRepository(db *gorm.DB) ProjectRepository {
	return &projectRepository{db: db}
}

// CreateProject creates a new project in the database.
func (r *projectRepository) CreateProject(project *models.Project) error {
	return r.db.Create(project).Error
}

// GetProjectByID retrieves a project by its ID, preloading Team.
func (r *projectRepository) GetProjectByID(id uint) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("Team").First(&project, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or return a specific "not found" error
		}
		return nil, err
	}
	return &project, nil
}

// GetProjectBySlug retrieves a project by its slug, preloading Team.
// Note: The current model has slug unique by team. This might need adjustment
// if slugs need to be globally unique or unique in a different context.
// For now, assuming slug is unique enough for this lookup, or the first found is acceptable.
func (r *projectRepository) GetProjectBySlug(slug string) (*models.Project, error) {
	var project models.Project
	err := r.db.Preload("Team").Where("slug = ?", slug).First(&project).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or return a specific "not found" error
		}
		return nil, err
	}
	return &project, nil
}

// GetProjectsByTeamID retrieves all projects for a given teamID.
func (r *projectRepository) GetProjectsByTeamID(teamID uint) ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Where("team_id = ?", teamID).Find(&projects).Error
	return projects, err
}

// UpdateProject updates an existing project in the database.
func (r *projectRepository) UpdateProject(project *models.Project) error {
	return r.db.Save(project).Error
}

// DeleteProject soft deletes a project by its ID.
func (r *projectRepository) DeleteProject(id uint) error {
	return r.db.Delete(&models.Project{}, id).Error
}

// GetAllProjects retrieves all projects from the database.
func (r *projectRepository) GetAllProjects() ([]models.Project, error) {
	var projects []models.Project
	err := r.db.Preload("Team").Find(&projects).Error // Also preloading Team for all projects
	return projects, err
}
