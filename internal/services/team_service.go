package services

import (
	"go-gin-gorm-api/internal/models"
	"go-gin-gorm-api/internal/repositories"
	"strings" // For basic slug generation
)

// TeamService defines the interface for team-related business logic.
type TeamService interface {
	CreateTeam(team *models.Team) error
	GetTeamByID(id uint) (*models.Team, error)
	GetTeamBySlug(slug string) (*models.Team, error)
	GetAllTeams() ([]models.Team, error)
	UpdateTeam(team *models.Team) error
	DeleteTeam(id uint) error
}

// teamService implements TeamService.
type teamService struct {
	teamRepo repositories.TeamRepository
}

// NewTeamService creates a new instance of TeamService.
func NewTeamService(repo repositories.TeamRepository) TeamService {
	return &teamService{teamRepo: repo}
}

// generateSlug creates a basic slug from a string.
// Replace with a more robust slugification library for production.
func generateSlug(name string) string {
	return strings.ToLower(strings.ReplaceAll(name, " ", "-"))
}

// CreateTeam creates a new team.
// If the team's slug is empty, it generates one from the name.
func (s *teamService) CreateTeam(team *models.Team) error {
	if team.Slug == "" && team.Name != "" {
		team.Slug = generateSlug(team.Name)
	}
	// Add any other business logic before creating, e.g., validation
	return s.teamRepo.CreateTeam(team)
}

// GetTeamByID retrieves a team by its ID.
func (s *teamService) GetTeamByID(id uint) (*models.Team, error) {
	return s.teamRepo.GetTeamByID(id)
}

// GetTeamBySlug retrieves a team by its slug.
func (s *teamService) GetTeamBySlug(slug string) (*models.Team, error) {
	return s.teamRepo.GetTeamBySlug(slug)
}

// GetAllTeams retrieves all teams.
func (s *teamService) GetAllTeams() ([]models.Team, error) {
	return s.teamRepo.GetAllTeams()
}

// UpdateTeam updates an existing team.
// Consider if slug should be updatable or if it needs regeneration.
func (s *teamService) UpdateTeam(team *models.Team) error {
	// Add any business logic before updating
	if team.Slug == "" && team.Name != "" { // Ensure slug is present if name is being updated
		team.Slug = generateSlug(team.Name)
	}
	return s.teamRepo.UpdateTeam(team)
}

// DeleteTeam deletes a team by its ID.
func (s *teamService) DeleteTeam(id uint) error {
	// Add any business logic before deleting (e.g., check for associated projects)
	return s.teamRepo.DeleteTeam(id)
}
