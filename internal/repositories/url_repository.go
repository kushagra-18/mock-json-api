package repositories

import (
	"errors"
	"go-gin-gorm-api/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// UrlRepository defines the interface for URL data operations.
type UrlRepository interface {
	CreateUrl(url *models.Url) error
	GetUrlByID(id uint) (*models.Url, error)
	GetUrlByPath(path string) (*models.Url, error) // Equivalent to 'findByUrl'
	GetUrlsByProjectID(projectID uint) ([]models.Url, error)
	UpdateUrl(url *models.Url) error
	DeleteUrl(id uint) error
	GetAllUrls() ([]models.Url, error)
	GetUrlsByTeamSlug(teamSlug string) ([]models.Url, error)
	GetUrlsByProjectSlug(projectSlug string) ([]models.Url, error)
	GetUrlByTeamSlugAndProjectSlugAndUrlPath(teamSlug, projectSlug, urlPath string) (*models.Url, error)
}

// urlRepository implements UrlRepository with GORM.
type urlRepository struct {
	db *gorm.DB
}

// NewUrlRepository creates a new instance of urlRepository.
func NewUrlRepository(db *gorm.DB) UrlRepository {
	return &urlRepository{db: db}
}

// CreateUrl creates a new URL in the database.
func (r *urlRepository) CreateUrl(url *models.Url) error {
	return r.db.Create(url).Error
}

// GetUrlByID retrieves a URL by its ID, preloading Project and MockContentList.
func (r *urlRepository) GetUrlByID(id uint) (*models.Url, error) {
	var url models.Url
	err := r.db.Preload(clause.Associations).Preload("MockContentList").First(&url, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or return a specific "not found" error
		}
		return nil, err
	}
	return &url, nil
}

// GetUrlByPath retrieves a URL by its path (URL field), preloading Project and MockContentList.
func (r *urlRepository) GetUrlByPath(path string) (*models.Url, error) {
	var url models.Url
	// Assuming 'URL' is the field name in the model struct
	err := r.db.Preload(clause.Associations).Preload("MockContentList").Where("url = ?", path).First(&url).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or return a specific "not found" error
		}
		return nil, err
	}
	return &url, nil
}

// GetUrlsByProjectID retrieves all URLs for a given projectID, preloading Project and MockContentList.
func (r *urlRepository) GetUrlsByProjectID(projectID uint) ([]models.Url, error) {
	var urls []models.Url
	err := r.db.Preload(clause.Associations).Preload("MockContentList").Where("project_id = ?", projectID).Find(&urls).Error
	return urls, err
}

// UpdateUrl updates an existing URL in the database.
func (r *urlRepository) UpdateUrl(url *models.Url) error {
	return r.db.Save(url).Error
}

// DeleteUrl soft deletes a URL by its ID.
func (r *urlRepository) DeleteUrl(id uint) error {
	// Also delete associated MockContent if desired (manual or via hooks/cascade)
	// For now, just deleting the URL itself.
	return r.db.Delete(&models.Url{}, id).Error
}

// GetAllUrls retrieves all URLs from the database, preloading Project and MockContentList.
func (r *urlRepository) GetAllUrls() ([]models.Url, error) {
	var urls []models.Url
	err := r.db.Preload(clause.Associations).Preload("MockContentList").Find(&urls).Error
	return urls, err
}

// GetUrlsByTeamSlug retrieves URLs by team slug.
// It joins with Project and Team, filters by team slug, and preloads associations.
func (r *urlRepository) GetUrlsByTeamSlug(teamSlug string) ([]models.Url, error) {
	var urls []models.Url
	// GORM will infer Project.Team from the model structure.
	// For Where conditions on joined tables, GORM often expects the actual table name or a pre-defined alias.
	// Assuming 'teams' is the table name for Team model.
	err := r.db.Joins("Project").Joins("Team", r.db.Where(&models.Team{Slug: teamSlug})).Preload("Project.Team").Preload("MockContentList").Find(&urls).Error
	// Alternative using explicit join condition if GORM has trouble inferring:
	// err := r.db.Joins("JOIN projects ON projects.id = urls.project_id").
	// Joins("JOIN teams ON teams.id = projects.team_id").
	// Where("teams.slug = ?", teamSlug).
	// Preload("Project.Team").Preload("MockContentList").Find(&urls).Error
	return urls, err
}

// GetUrlsByProjectSlug retrieves URLs by project slug.
// It joins with Project, filters by project slug, and preloads associations.
func (r *urlRepository) GetUrlsByProjectSlug(projectSlug string) ([]models.Url, error) {
	var urls []models.Url
	// Assuming 'projects' is the table name for Project model.
	err := r.db.Joins("Project", r.db.Where(&models.Project{Slug: projectSlug})).Preload("Project").Preload("MockContentList").Find(&urls).Error
	// Alternative explicit join:
	// err := r.db.Joins("JOIN projects ON projects.id = urls.project_id").
	// Where("projects.slug = ?", projectSlug).
	// Preload("Project").Preload("MockContentList").Find(&urls).Error
	return urls, err
}

// GetUrlByTeamSlugAndProjectSlugAndUrlPath retrieves a single URL by team slug, project slug, and URL path.
// It joins with Project and Team, filters by these slugs and path, and preloads associations.
func (r *urlRepository) GetUrlByTeamSlugAndProjectSlugAndUrlPath(teamSlug, projectSlug, urlPath string) (*models.Url, error) {
	var url models.Url
	// Using actual table names in Where for clarity and robustness with Joins.
	// GORM's default table name for Project model is 'projects', for Team model is 'teams'.
	// The 'urls.url' refers to the 'url' column in the 'urls' table.
	err := r.db.Joins("JOIN projects ON projects.id = urls.project_id").
		Joins("JOIN teams ON teams.id = projects.team_id").
		Where("teams.slug = ? AND projects.slug = ? AND urls.url = ?", teamSlug, projectSlug, urlPath).
		Preload("Project.Team").Preload("MockContentList").First(&url).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or return a specific "not found" error
		}
		return nil, err
	}
	return &url, nil
}
