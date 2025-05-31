package services

// MockFakerService is a manual mock for FakerService.
type MockFakerService struct {
	ProcessDSLFunc func(dslString string) (string, error)
}

// ProcessDSL delegates to ProcessDSLFunc if it's set.
// Otherwise, it might panic or return a default value, depending on test needs.
// For robustness in tests, it's often better if the test explicitly sets ProcessDSLFunc.
func (m *MockFakerService) ProcessDSL(dslString string) (string, error) {
	if m.ProcessDSLFunc != nil {
		return m.ProcessDSLFunc(dslString)
	}
	// Consider returning a default error or specific mock behavior if not set,
	// e.g., return "", fmt.Errorf("ProcessDSLFunc not implemented in mock")
	// For now, panicking makes it clear that the test setup is incomplete.
	panic("MockFakerService.ProcessDSLFunc is not set")
}
