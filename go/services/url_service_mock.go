package services

import "mockapi/models"

// MockURLService is a manual mock for URLService.
type MockURLService struct {
	FindByProjectIDAndURLFunc        func(projectID uint, urlPath string) (*models.Url, error)
	CreateURLFunc                    func(url *models.Url, projectID uint) error
	GetURLByIDFunc                   func(id uint) (*models.Url, error)
	GetURLByTeamSlugProjectSlugAndPathFunc func(teamSlug, projectSlug, path string) (*models.Url, error)
	IncrementRequestStatsFunc        func(urlID uint) error
	// Add other methods used by MockContentController if any
}

func (m *MockURLService) FindByProjectIDAndURL(projectID uint, urlPath string) (*models.Url, error) {
	if m.FindByProjectIDAndURLFunc != nil {
		return m.FindByProjectIDAndURLFunc(projectID, urlPath)
	}
	panic("MockURLService.FindByProjectIDAndURLFunc is not set")
}

func (m *MockURLService) CreateURL(url *models.Url, projectID uint) error {
	if m.CreateURLFunc != nil {
		return m.CreateURLFunc(url, projectID)
	}
	panic("MockURLService.CreateURLFunc is not set")
}

func (m *MockURLService) GetURLByID(id uint) (*models.Url, error) {
	if m.GetURLByIDFunc != nil {
		return m.GetURLByIDFunc(id)
	}
	panic("MockURLService.GetURLByIDFunc is not set")
}

func (m *MockURLService) GetURLByTeamSlugProjectSlugAndPath(teamSlug, projectSlug, path string) (*models.Url, error) {
	if m.GetURLByTeamSlugProjectSlugAndPathFunc != nil {
		return m.GetURLByTeamSlugProjectSlugAndPathFunc(teamSlug, projectSlug, path)
	}
	panic("MockURLService.GetURLByTeamSlugProjectSlugAndPathFunc is not set")
}

func (m *MockURLService) IncrementRequestStats(urlID uint) error {
	if m.IncrementRequestStatsFunc != nil {
		return m.IncrementRequestStatsFunc(urlID)
	}
	panic("MockURLService.IncrementRequestStatsFunc is not set")
}

// Ensure this mock implements all methods of URLService that are actually called by the controller.
// SaveMockContent uses: FindByProjectIDAndURL, CreateURL
// UpdateMockContent uses: GetURLByID
// GetMockedJSON uses: GetURLByTeamSlugProjectSlugAndPath, IncrementRequestStats
// The mock includes these. Add others if controller logic expands.
