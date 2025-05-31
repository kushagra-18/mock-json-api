package services

// MockRedisService is a manual mock for RedisService.
type MockRedisService struct {
	RateLimitFunc      func(key string, limit int, window int64) (bool, error)
	CreateRedisKeyFunc func(args ...string) string
	// Add other methods used by MockContentController if any (e.g., related to caching URL data if implemented)
}

func (m *MockRedisService) RateLimit(key string, limit int, window int64) (bool, error) {
	if m.RateLimitFunc != nil {
		return m.RateLimitFunc(key, limit, window)
	}
	panic("MockRedisService.RateLimitFunc is not set")
}

func (m *MockRedisService) CreateRedisKey(args ...string) string {
	if m.CreateRedisKeyFunc != nil {
		return m.CreateRedisKeyFunc(args...)
	}
	panic("MockRedisService.CreateRedisKeyFunc is not set")
}

// Ensure this mock implements all methods of RedisService that are actually called by the controller.
// GetMockedJSON uses: RateLimit, CreateRedisKey
// SaveMockContent and UpdateMockContent do not directly call RedisService methods in the provided code.
// The mock includes these. Add others if controller logic expands.
