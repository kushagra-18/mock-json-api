package services

import "mockapi/models"

// MockProjectService is a manual mock for ProjectService.
type MockProjectService struct {
	GetProjectBySlugFunc                 func(slug string) (*models.Project, error)
	GetProjectByTeamSlugAndProjectSlugFunc func(teamSlug, projectSlug string) (*models.Project, error)
	// Add other methods used by MockContentController if any
}

func (m *MockProjectService) GetProjectBySlug(slug string) (*models.Project, error) {
	if m.GetProjectBySlugFunc != nil {
		return m.GetProjectBySlugFunc(slug)
	}
	panic("MockProjectService.GetProjectBySlugFunc is not set")
}

func (m *MockProjectService) GetProjectByTeamSlugAndProjectSlug(teamSlug, projectSlug string) (*models.Project, error) {
	if m.GetProjectByTeamSlugAndProjectSlugFunc != nil {
		return m.GetProjectByTeamSlugAndProjectSlugFunc(teamSlug, projectSlug)
	}
	panic("MockProjectService.GetProjectByTeamSlugAndProjectSlugFunc is not set")
}
// Ensure this mock implements all methods of ProjectService that are actually called by the controller.
// For instance, if ProjectService has methods like `GetProjectByID`, `CreateProject`, etc.,
// and if those are potentially called through the controller's methods being tested,
// they would need to be added here. For SaveMockContent and UpdateMockContent,
// GetProjectBySlug is the primary one. GetMockedJSON uses GetProjectByTeamSlugAndProjectSlug.
// For now, this is a minimal mock based on controller usage in Save/Update.
