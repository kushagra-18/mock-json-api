package services

import "mockapi/models"

// MockRequestLogService is a manual mock for RequestLogService.
type MockRequestLogService struct {
	SaveRequestLogFunc func(logEntry *models.RequestLog) error
	// Add other methods used by MockContentController if any
}

func (m *MockRequestLogService) SaveRequestLog(logEntry *models.RequestLog) error {
	if m.SaveRequestLogFunc != nil {
		return m.SaveRequestLogFunc(logEntry)
	}
	panic("MockRequestLogService.SaveRequestLogFunc is not set")
}

// Ensure this mock implements all methods of RequestLogService that are actually called by the controller.
// GetMockedJSON uses: SaveRequestLog (via finalizeRequestLog)
// SaveMockContent and UpdateMockContent do not directly call RequestLogService methods in the provided code,
// but they might if extensive logging/auditing were added there.
// The mock includes SaveRequestLog. Add others if controller logic expands.
