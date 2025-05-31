package services

import (
	"fmt"

	"gorm.io/gorm"
	"mockapi/models" // Assuming module name is mockapi
)

// ProxyService handles business logic related to forward proxies.
type ProxyService struct {
	DB *gorm.DB
}

// NewProxyService creates a new ProxyService.
func NewProxyService(db *gorm.DB) *ProxyService {
	return &ProxyService{DB: db}
}

// GetForwardProxyByProjectID retrieves the forward proxy settings for a given project ID.
func (s *ProxyService) GetForwardProxyByProjectID(projectID uint) (*models.ForwardProxy, error) {
	var proxy models.ForwardProxy
	// Assuming a project has one ForwardProxy. If it can exist or not, FirstOrInit or FirstOrCreate might be options
	// or just First and check for ErrRecordNotFound.
	if err := s.DB.Where("project_id = ?", projectID).First(&proxy).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// It's okay for a project not to have a forward proxy configured.
			return nil, nil // Or return a specific error like ErrForwardProxyNotConfigured
		}
		return nil, fmt.Errorf("failed to retrieve forward proxy for project ID %d: %w", projectID, err)
	}
	return &proxy, nil
}

// CreateForwardProxy creates or updates the forward proxy settings for a project.
// Since a project has at most one forward proxy (due to unique project_id), this can be an upsert.
func (s *ProxyService) CreateForwardProxy(proxy *models.ForwardProxy, projectID uint) (*models.ForwardProxy, error) {
	if proxy == nil {
		return nil, fmt.Errorf("proxy data cannot be nil")
	}
	proxy.ProjectID = projectID

	// Check if a proxy already exists for this project
	var existingProxy models.ForwardProxy
	err := s.DB.Where("project_id = ?", projectID).First(&existingProxy).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("error checking for existing forward proxy for project ID %d: %w", projectID, err)
	}

	if err == gorm.ErrRecordNotFound { // No existing proxy, create new
		if createErr := s.DB.Create(proxy).Error; createErr != nil {
			return nil, fmt.Errorf("failed to create forward proxy for project ID %d: %w", projectID, createErr)
		}
		return proxy, nil
	}

	// Existing proxy found, update it
	existingProxy.Domain = proxy.Domain // Update relevant fields
	// Update other fields from 'proxy' to 'existingProxy' as necessary

	if updateErr := s.DB.Save(&existingProxy).Error; updateErr != nil {
		return nil, fmt.Errorf("failed to update forward proxy for project ID %d: %w", projectID, updateErr)
	}
	return &existingProxy, nil
}

// UpdateForwardProxy updates an existing forward proxy.
// This is more generic if more fields were part of ForwardProxy.
// Given ForwardProxy only has Domain and ProjectID (which shouldn't change post-creation this way),
// CreateForwardProxy effectively handles updates for the 'Domain'.
// This method can be used if a direct update by proxy ID is preferred.
func (s *ProxyService) UpdateForwardProxy(proxyToUpdate *models.ForwardProxy) (*models.ForwardProxy, error) {
	if proxyToUpdate == nil || proxyToUpdate.ID == 0 {
		return nil, fmt.Errorf("forward proxy data is invalid or ID is missing for update")
	}

	// Optional: Check if ProjectID is being changed, which might be disallowed.
	// var originalProxy models.ForwardProxy
	// if err := s.DB.First(&originalProxy, proxyToUpdate.ID).Error; err != nil {
	// 	 return nil, fmt.Errorf("original forward proxy with ID %d not found: %w", proxyToUpdate.ID, err)
	// }
	// if originalProxy.ProjectID != proxyToUpdate.ProjectID {
	//    return nil, fmt.Errorf("projectID of a forward proxy cannot be changed during update")
	// }


	if err := s.DB.Save(proxyToUpdate).Error; err != nil {
		return nil, fmt.Errorf("failed to update forward proxy with ID %d: %w", proxyToUpdate.ID, err)
	}
	return proxyToUpdate, nil
}

// DeleteForwardProxyByProjectID deletes the forward proxy for a given project ID.
func (s *ProxyService) DeleteForwardProxyByProjectID(projectID uint) error {
	result := s.DB.Where("project_id = ?", projectID).Delete(&models.ForwardProxy{})
	if result.Error != nil {
		return fmt.Errorf("failed to delete forward proxy for project ID %d: %w", projectID, result.Error)
	}
	// No error if RowsAffected is 0, as it means no proxy was configured.
	return nil
}

// DeleteForwardProxyByID deletes a forward proxy by its own ID.
func (s *ProxyService) DeleteForwardProxyByID(proxyID uint) error {
    result := s.DB.Delete(&models.ForwardProxy{}, proxyID)
    if result.Error != nil {
        return fmt.Errorf("failed to delete forward proxy with ID %d: %w", proxyID, result.Error)
    }
    if result.RowsAffected == 0 {
        return fmt.Errorf("forward proxy with ID %d not found for deletion: %w", proxyID, gorm.ErrRecordNotFound)
    }
    return nil
}
