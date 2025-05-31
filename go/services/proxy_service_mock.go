package services

import "mockapi/models"

// MockProxyService is a manual mock for ProxyService.
type MockProxyService struct {
	GetForwardProxyByProjectIDFunc func(projectID uint) (*models.ForwardProxy, error)
	// Add other methods used by MockContentController if any
}

func (m *MockProxyService) GetForwardProxyByProjectID(projectID uint) (*models.ForwardProxy, error) {
	if m.GetForwardProxyByProjectIDFunc != nil {
		return m.GetForwardProxyByProjectIDFunc(projectID)
	}
	panic("MockProxyService.GetForwardProxyByProjectIDFunc is not set")
}

// Ensure this mock implements all methods of ProxyService that are actually called by the controller.
// GetMockedJSON uses: GetForwardProxyByProjectID
// SaveMockContent and UpdateMockContent do not directly call ProxyService methods in the provided code.
// The mock includes this. Add others if controller logic expands.
