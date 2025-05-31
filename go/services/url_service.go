package services

import (
	"database/sql"
	"fmt"
	"strings"
	// "time" // No longer used directly

	"gorm.io/gorm"
	"mockapi/dtos"
	"mockapi/models" // Assuming module name is mockapi
)

// URLService handles business logic related to URLs.
type URLService struct {
	DB           *gorm.DB
	RedisService *RedisService
}

// NewURLService creates a new URLService.
func NewURLService(db *gorm.DB, redisService *RedisService) *URLService {
	return &URLService{DB: db, RedisService: redisService}
}

// GetURLByID retrieves a URL by its ID.
func (s *URLService) GetURLByID(id uint) (*models.Url, error) {
	var url models.Url
	if err := s.DB.Preload("MockContents").First(&url, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("url with ID %d not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to retrieve url with ID %d: %w", id, err)
	}
	return &url, nil
}

// GetURLByProjectSlugAndPath retrieves a URL by its project's slug and its path.
// The path should be the full path including query parameters if they are part of the URL key.
func (s *URLService) GetURLByProjectSlugAndPath(projectSlug, path string) (*models.Url, error) {
	var url models.Url
	err := s.DB.Joins("JOIN projects ON projects.id = urls.project_id").
		Where("projects.slug = ? AND urls.url = ?", projectSlug, path).
		Preload("MockContents"). // Preload mock contents
		First(&url).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("url with path '%s' under project '%s' not found: %w", path, projectSlug, err)
		}
		return nil, fmt.Errorf("failed to retrieve url with path '%s' under project '%s': %w", path, projectSlug, err)
	}
	return &url, nil
}

// GetURLByTeamSlugProjectSlugAndPath retrieves a URL by team slug, project slug, and URL path.
func (s *URLService) GetURLByTeamSlugProjectSlugAndPath(teamSlug, projectSlug, path string) (*models.Url, error) {
	var url models.Url
	err := s.DB.Joins("JOIN projects ON projects.id = urls.project_id").
		Joins("JOIN teams ON teams.id = projects.team_id").
		Where("teams.slug = ? AND projects.slug = ? AND urls.url = ?", teamSlug, projectSlug, path).
		Preload("MockContents").
		First(&url).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("url with path '%s' for project '%s' (team '%s') not found: %w", path, projectSlug, teamSlug, err)
		}
		return nil, fmt.Errorf("failed to retrieve url with path '%s' for project '%s' (team '%s'): %w", path, projectSlug, teamSlug, err)
	}
	return &url, nil
}

// CreateURL creates a new URL for a given project.
func (s *URLService) CreateURL(url *models.Url, projectID uint) error {
	if url == nil {
		return fmt.Errorf("url data cannot be nil")
	}
	url.ProjectID = projectID
	if url.Status == "" {
		url.Status = models.StatusOK // Default status
	}

	if err := s.DB.Create(url).Error; err != nil {
		// Consider checking for unique constraint violation errors specifically
		return fmt.Errorf("failed to create url: %w", err)
	}
	return nil
}

// UpdateURL updates an existing URL using data from URLDataDTO.
func (s *URLService) UpdateURL(urlID uint, dto dtos.URLDataDTO) (*models.Url, error) {
	var urlToUpdate models.Url
	if err := s.DB.First(&urlToUpdate, urlID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("url with ID %d not found for update: %w", urlID, err)
		}
		return nil, fmt.Errorf("failed to find url with ID %d for update: %w", urlID, err)
	}

	if dto.Description != nil {
		urlToUpdate.Description = *dto.Description
	}
	if dto.Name != nil {
		urlToUpdate.Name = *dto.Name
	}
	if dto.Requests != nil {
		urlToUpdate.Requests = sql.NullInt64{Int64: *dto.Requests, Valid: true}
	}
	if dto.Time != nil {
		urlToUpdate.Time = sql.NullInt64{Int64: *dto.Time, Valid: true}
	}
	if dto.Status != nil {
		// Basic validation, can be expanded
		normalizedStatus := strings.ToUpper(*dto.Status)
		// This is a simplified check. Ideally, validate against the defined StatusCode constants.
		// For now, just ensuring it's not empty.
		if normalizedStatus != "" {
			urlToUpdate.Status = models.StatusCode(normalizedStatus)
		} else {
            // Keep existing status or handle as an error if status update is invalid
        }
	}


	if err := s.DB.Save(&urlToUpdate).Error; err != nil {
		return nil, fmt.Errorf("failed to update url with ID %d: %w", urlID, err)
	}
	return &urlToUpdate, nil
}

// IsRateLimited checks if a request for a given IP and path is rate-limited.
// It uses the RedisService for the underlying rate-limiting logic.
func (s *URLService) IsRateLimited(ip, path string, allowedRequests int, timeWindowSeconds int64) (bool, error) {
	if s.RedisService == nil {
		return false, fmt.Errorf("redis service is not initialized")
	}
	// Normalize path to be part of the key, e.g., remove leading/trailing slashes or case transform
	normalizedPath := strings.ToLower(strings.Trim(path, "/"))
	key := s.RedisService.CreateRedisKey("ratelimit", ip, normalizedPath)

	limited, err := s.RedisService.RateLimit(key, allowedRequests, timeWindowSeconds)
	if err != nil {
		return true, fmt.Errorf("error checking rate limit: %w", err) // Fail closed
	}
	return limited, nil
}

// IncrementRequestStats increments the request count and updates last accessed time for a URL.
func (s *URLService) IncrementRequestStats(urlID uint) error {
	// Using .Updates to only update specified fields and trigger hooks if necessary
	// It's generally safer and more explicit than .Save for partial updates.
	// GORM handles nullable types like sql.NullInt64 correctly with Updates.
	// Here we assume 'requests' should be incremented.
	// 'time' might represent 'last_accessed_time' or 'average_response_time'.
	// If 'time' is 'last_accessed_time', it should be set to current time.
	// For this example, let's assume 'requests' is incremented.
	// The 'Time' field was sql.NullInt64 in the model, its meaning here is a bit ambiguous from the DTO.
	// If it was 'average response time', it would be calculated differently.
	// If it was 'last request processing time', it would be set per request.
	// For now, just incrementing 'Requests'.
	result := s.DB.Model(&models.Url{}).Where("id = ?", urlID).UpdateColumn("requests", gorm.Expr("requests + 1"))
	// To update 'UpdatedAt' timestamp as well, use .Updates instead of .UpdateColumn
	// result := s.DB.Model(&models.Url{}).Where("id = ?", urlID).Updates(map[string]interface{}{
	// 	"requests": gorm.Expr("requests + 1"),
	// 	"updated_at": time.Now(), // Explicitly set if not relying on GORM's auto-update for this specific action
	// })


	if result.Error != nil {
		return fmt.Errorf("failed to increment request count for url ID %d: %w", urlID, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("url with ID %d not found for incrementing request stats: %w", urlID, gorm.ErrRecordNotFound)
	}
	return nil
}

// FindByProjectIDAndURL finds a URL by its project ID and the exact URL string.
func (s *URLService) FindByProjectIDAndURL(projectID uint, urlPath string) (*models.Url, error) {
    var url models.Url
    err := s.DB.Where("project_id = ? AND url = ?", projectID, urlPath).
        Preload("MockContents").
        First(&url).Error

    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, fmt.Errorf("url with path '%s' under project ID %d not found: %w", urlPath, projectID, err)
        }
        return nil, fmt.Errorf("failed to retrieve url with path '%s' under project ID %d: %w", urlPath, projectID, err)
    }
    return &url, nil
}

// DeleteURL deletes a URL by its ID. It will also need to handle related MockContents.
func (s *URLService) DeleteURL(urlID uint) error {
    // GORM's default behavior with associations might require explicit deletion of MockContents
    // or a database-level CASCADE constraint. For now, let's assume GORM handles it or
    // it's handled by cascading deletes in the DB.
    // If not, one would first delete MockContents:
    // s.DB.Where("url_id = ?", urlID).Delete(&models.MockContent{})

    result := s.DB.Delete(&models.Url{}, urlID)
    if result.Error != nil {
        return fmt.Errorf("failed to delete url with ID %d: %w", urlID, result.Error)
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("url with ID %d not found for deletion: %w", urlID, gorm.ErrRecordNotFound)
    }
    return nil
}

// GetURLsByProjectID retrieves all URLs for a given project ID.
func (s *URLService) GetURLsByProjectID(projectID uint) ([]models.Url, error) {
    var urls []models.Url
    err := s.DB.Where("project_id = ?", projectID).Find(&urls).Error
    if err != nil {
        return nil, fmt.Errorf("failed to retrieve urls for project ID %d: %w", projectID, err)
    }
    return urls, nil
}
