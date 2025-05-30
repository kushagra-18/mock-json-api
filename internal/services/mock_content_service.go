package services

import (
	"fmt"
	"go-gin-gorm-api/internal/models"
	"go-gin-gorm-api/internal/repositories"
	"math/rand"
	"time"
)

// MockContentService defines the interface for mock content-related business logic.
type MockContentService interface {
	CreateMockContent(mockContent *models.MockContent, urlID uint) (*models.MockContent, error)
	GetMockContentByID(id uint) (*models.MockContent, error)
	GetMockContentsByUrlID(urlID uint) ([]models.MockContent, error)
	UpdateMockContent(mockContent *models.MockContent) (*models.MockContent, error)
	DeleteMockContent(id uint) error
	SimulateLatency(mockContent *models.MockContent)
	SelectRandomMockContent(mockContents []models.MockContent) *models.MockContent
}

// mockContentService implements MockContentService.
type mockContentService struct {
	mockRepo repositories.MockContentRepository
	urlRepo  repositories.UrlRepository // For potential validation or URL-related logic
}

// NewMockContentService creates a new instance of MockContentService.
func NewMockContentService(mockRepo repositories.MockContentRepository, urlRepo repositories.UrlRepository) MockContentService {
	return &mockContentService{mockRepo: mockRepo, urlRepo: urlRepo}
}

// CreateMockContent creates new mock content, associating it with a URL.
func (s *mockContentService) CreateMockContent(mockContent *models.MockContent, urlID uint) (*models.MockContent, error) {
	// Optional: Validate urlID exists
	// url, err := s.urlRepo.GetUrlByID(urlID)
	// if err != nil {
	// 	return nil, fmt.Errorf("error validating URL ID: %w", err)
	// }
	// if url == nil {
	// 	return nil, fmt.Errorf("URL with ID %d not found", urlID)
	// }

	mockContent.URLID = urlID
	// Timestamps are usually handled by GORM's default behavior if the fields are time.Time
	// For explicit control:
	// now := time.Now()
	// mockContent.CreatedAt = now
	// mockContent.UpdatedAt = now

	err := s.mockRepo.CreateMockContent(mockContent)
	if err != nil {
		return nil, err
	}
	return mockContent, nil
}

// GetMockContentByID retrieves mock content by its ID.
func (s *mockContentService) GetMockContentByID(id uint) (*models.MockContent, error) {
	return s.mockRepo.GetMockContentByID(id)
}

// GetMockContentsByUrlID retrieves mock contents by URL ID.
func (s *mockContentService) GetMockContentsByUrlID(urlID uint) ([]models.MockContent, error) {
	return s.mockRepo.GetMockContentsByUrlID(urlID)
}

// UpdateMockContent updates existing mock content.
func (s *mockContentService) UpdateMockContent(mockContent *models.MockContent) (*models.MockContent, error) {
	// Timestamps are usually handled by GORM's default behavior if the fields are time.Time
	// For explicit control:
	// mockContent.UpdatedAt = time.Now()
	err := s.mockRepo.UpdateMockContent(mockContent)
	if err != nil {
		return nil, err
	}
	return mockContent, nil
}

// DeleteMockContent deletes mock content by its ID.
func (s *mockContentService) DeleteMockContent(id uint) error {
	return s.mockRepo.DeleteMockContent(id)
}

// SimulateLatency introduces a delay based on the mock content's latency configuration.
func (s *mockContentService) SimulateLatency(mockContent *models.MockContent) {
	if mockContent == nil || mockContent.Latency <= 0 {
		return
	}
	time.Sleep(time.Duration(mockContent.Latency) * time.Millisecond)
}

// SelectRandomMockContent selects a mock content item based on weighted randomness.
func (s *mockContentService) SelectRandomMockContent(mockContents []models.MockContent) *models.MockContent {
	if len(mockContents) == 0 {
		return nil
	}

	totalWeight := 0
	for _, mc := range mockContents {
		if mc.Randomness > 0 { // Consider only positive randomness values
			totalWeight += mc.Randomness
		}
	}

	if totalWeight <= 0 {
		// If no items have positive randomness, or list is all zero/negative randomness,
		// return the first item or nil (or implement other default behavior).
		// For now, returning the first item with non-zero ID if available, else nil.
		for _, mc := range mockContents {
			if mc.ID != 0 { // Check if it's a valid item
				return &mc
			}
		}
		return nil
	}

	// Seed the random number generator.
	// For better randomness, seed should ideally be done once at application start.
	// However, for simplicity in this context, seeding here.
	// rand.Seed(time.Now().UnixNano()) // Deprecated in Go 1.20+
	// In Go 1.20+, rand.New(rand.NewSource(time.Now().UnixNano())) is preferred for local instances
	// or just use global rand which is auto-seeded.

	r := rand.Intn(totalWeight) // Generates a random number in [0, totalWeight)

	for _, mc := range mockContents {
		if mc.Randomness > 0 {
			r -= mc.Randomness
			if r < 0 {
				return &mc
			}
		}
	}

	// Fallback in case of rounding errors or unexpected behavior,
	// though theoretically, one item should be selected if totalWeight > 0.
	// Return the last item with positive randomness if loop finishes.
	for i := len(mockContents) - 1; i >= 0; i-- {
		if mockContents[i].Randomness > 0 {
			return &mockContents[i]
		}
	}
    // If all randomness are zero or negative (already handled by totalWeight <=0), or empty list
	if len(mockContents) > 0 {
        return &mockContents[0] // return first one as default
    }

	return nil
}
