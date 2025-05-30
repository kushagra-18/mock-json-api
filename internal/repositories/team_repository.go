package repositories

import (
	"errors"
	"go-gin-gorm-api/internal/models"

	"gorm.io/gorm"
)

// TeamRepository defines the interface for team data operations.
type TeamRepository interface {
	CreateTeam(team *models.Team) error
	GetTeamByID(id uint) (*models.Team, error)
	GetTeamBySlug(slug string) (*models.Team, error)
	UpdateTeam(team *models.Team) error
	DeleteTeam(id uint) error
	GetAllTeams() ([]models.Team, error)
}

// teamRepository implements TeamRepository with GORM.
type teamRepository struct {
	db *gorm.DB
}

// NewTeamRepository creates a new instance of teamRepository.
func NewTeamRepository(db *gorm.DB) TeamRepository {
	return &teamRepository{db: db}
}

// CreateTeam creates a new team in the database.
func (r *teamRepository) CreateTeam(team *models.Team) error {
	return r.db.Create(team).Error
}

// GetTeamByID retrieves a team by its ID, preloading Projects.
func (r *teamRepository) GetTeamByID(id uint) (*models.Team, error) {
	var team models.Team
	err := r.db.Preload("Projects").First(&team, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or return a specific "not found" error
		}
		return nil, err
	}
	return &team, nil
}

// GetTeamBySlug retrieves a team by its slug, preloading Projects.
func (r *teamRepository) GetTeamBySlug(slug string) (*models.Team, error) {
	var team models.Team
	err := r.db.Preload("Projects").Where("slug = ?", slug).First(&team).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or return a specific "not found" error
		}
		return nil, err
	}
	return &team, nil
}

// UpdateTeam updates an existing team in the database.
func (r *teamRepository) UpdateTeam(team *models.Team) error {
	return r.db.Save(team).Error
}

// DeleteTeam soft deletes a team by its ID.
func (r *teamRepository) DeleteTeam(id uint) error {
	return r.db.Delete(&models.Team{}, id).Error
}

// GetAllTeams retrieves all teams from the database.
func (r *teamRepository) GetAllTeams() ([]models.Team, error) {
	var teams []models.Team
	err := r.db.Find(&teams).Error
	return teams, err
}
