package services

import (
	"fmt"
	"go-gin-gorm-api/internal/models"
	"go-gin-gorm-api/internal/repositories"
	"net/http"
	"strings"
)

// UrlService defines the interface for URL-related business logic.
type UrlService interface {
	GetUrlByDetails(teamSlug, projectSlug, urlPath string) (*models.Url, error)
	GetUrlByPath(path string) (*models.Url, error)
	ExtractFullUrl(r *http.Request) string
	CreateUrl(url *models.Url, projectID uint) (*models.Url, error)
	UpdateUrl(url *models.Url) (*models.Url, error)
	DeleteUrl(id uint) error
	GetUrlsByProjectID(projectID uint) ([]models.Url, error)
}

// urlService implements UrlService.
type urlService struct {
	urlRepo     repositories.UrlRepository
	projectRepo repositories.ProjectRepository // For potential validation or project-related logic
}

// NewUrlService creates a new instance of UrlService.
func NewUrlService(urlRepo repositories.UrlRepository, projectRepo repositories.ProjectRepository) UrlService {
	return &urlService{urlRepo: urlRepo, projectRepo: projectRepo}
}

// GetUrlByDetails retrieves a URL based on team slug, project slug, and URL path.
func (s *urlService) GetUrlByDetails(teamSlug, projectSlug, urlPath string) (*models.Url, error) {
	return s.urlRepo.GetUrlByTeamSlugAndProjectSlugAndUrlPath(teamSlug, projectSlug, urlPath)
}

// GetUrlByPath retrieves a URL by its path.
func (s *urlService) GetUrlByPath(path string) (*models.Url, error) {
	return s.urlRepo.GetUrlByPath(path)
}

// ExtractFullUrl constructs the full path including query string from an HTTP request.
func (s *urlService) ExtractFullUrl(r *http.Request) string {
	path := r.URL.Path
	if r.URL.RawQuery != "" {
		return fmt.Sprintf("%s?%s", path, r.URL.RawQuery)
	}
	// Normalize: ensure leading slash if path is not empty
	if path != "" && !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return path
}

// CreateUrl creates a new URL, associating it with a project.
func (s *urlService) CreateUrl(url *models.Url, projectID uint) (*models.Url, error) {
	// Optional: Validate projectID exists
	// project, err := s.projectRepo.GetProjectByID(projectID)
	// if err != nil {
	// return nil, fmt.Errorf("error validating project ID: %w", err)
	// }
	// if project == nil {
	// return nil, fmt.Errorf("project with ID %d not found", projectID)
	// }

	url.ProjectID = projectID
	err := s.urlRepo.CreateUrl(url)
	if err != nil {
		return nil, err
	}
	return url, nil
}

// UpdateUrl updates an existing URL.
func (s *urlService) UpdateUrl(url *models.Url) (*models.Url, error) {
	err := s.urlRepo.UpdateUrl(url)
	if err != nil {
		return nil, err
	}
	return url, nil
}

// DeleteUrl deletes a URL by its ID.
func (s *urlService) DeleteUrl(id uint) error {
	return s.urlRepo.DeleteUrl(id)
}

// GetUrlsByProjectID retrieves URLs by project ID.
func (s *urlService) GetUrlsByProjectID(projectID uint) ([]models.Url, error) {
	return s.urlRepo.GetUrlsByProjectID(projectID)
}
