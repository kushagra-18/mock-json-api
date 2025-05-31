package services

import "mockapi/models"

// MockMockContentService is a manual mock for MockContentService.
type MockMockContentService struct {
	SaveMockContentListFunc   func(mockContents []models.MockContent, urlID uint) ([]models.MockContent, error)
	UpdateMockContentListFunc func(mockContents []models.MockContent, urlID uint) ([]models.MockContent, error)
	SelectRandomMockContentFunc func(contents []models.MockContent) *models.MockContent
	SimulateLatencyFunc         func(latency int64)
	// Add other methods used by MockContentController if any
}

func (m *MockMockContentService) SaveMockContentList(mockContents []models.MockContent, urlID uint) ([]models.MockContent, error) {
	if m.SaveMockContentListFunc != nil {
		return m.SaveMockContentListFunc(mockContents, urlID)
	}
	panic("MockMockContentService.SaveMockContentListFunc is not set")
}

func (m *MockMockContentService) UpdateMockContentList(mockContents []models.MockContent, urlID uint) ([]models.MockContent, error) {
	if m.UpdateMockContentListFunc != nil {
		return m.UpdateMockContentListFunc(mockContents, urlID)
	}
	panic("MockMockContentService.UpdateMockContentListFunc is not set")
}

func (m *MockMockContentService) SelectRandomMockContent(contents []models.MockContent) *models.MockContent {
	if m.SelectRandomMockContentFunc != nil {
		return m.SelectRandomMockContentFunc(contents)
	}
	panic("MockMockContentService.SelectRandomMockContentFunc is not set")
}

func (m *MockMockContentService) SimulateLatency(latency int64) {
	if m.SimulateLatencyFunc != nil {
		m.SimulateLatencyFunc(latency)
		return
	}
	// Default behavior can be a no-op for latency simulation in many tests
	// or panic if strict control is needed: panic("MockMockContentService.SimulateLatencyFunc is not set")
}

// Ensure this mock implements all methods of MockContentService that are actually called by the controller.
// SaveMockContent uses: SaveMockContentList
// UpdateMockContent uses: UpdateMockContentList
// GetMockedJSON uses: SelectRandomMockContent, SimulateLatency
// The mock includes these. Add others if controller logic expands.
