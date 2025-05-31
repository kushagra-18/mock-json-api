package services

import (
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
	"mockapi/models" // Assuming module name is mockapi
)

// MockContentService handles business logic related to mock contents.
type MockContentService struct {
	DB *gorm.DB
}

// NewMockContentService creates a new MockContentService.
func NewMockContentService(db *gorm.DB) *MockContentService {
	return &MockContentService{DB: db}
}

// SaveMockContentList saves a list of new mock contents for a given URL.
// Assumes these are all new entries.
func (s *MockContentService) SaveMockContentList(mockContents []models.MockContent, urlID uint) ([]models.MockContent, error) {
	if len(mockContents) == 0 {
		return []models.MockContent{}, nil
	}

	for i := range mockContents {
		mockContents[i].UrlID = urlID
		// Ensure ID is zero for GORM to treat as new record if it's accidentally set
		mockContents[i].ID = 0
		mockContents[i].BaseModel.ID = 0
	}

	// Using CreateInBatches can be efficient for large lists.
	// However, for simplicity and returning created objects with IDs, iterating or single Create might be okay.
	// GORM's Create can handle a slice directly.
	if err := s.DB.Create(&mockContents).Error; err != nil {
		return nil, fmt.Errorf("failed to save mock content list for url ID %d: %w", urlID, err)
	}
	return mockContents, nil
}

// UpdateMockContentList updates a list of mock contents for a given URL.
// This function needs to decide whether to create new ones, update existing ones, or delete old ones.
// A common strategy: delete all existing for the URL, then create new ones.
// Or, match by ID if provided, create if ID is zero, update if ID exists.
// For this implementation, let's go with a "delete all then create all" strategy for simplicity,
// which matches some of the behavior in the original Java `MockContentService.save(List<MockContent>, Url)`.
func (s *MockContentService) UpdateMockContentList(mockContents []models.MockContent, urlID uint) ([]models.MockContent, error) {
	// Start a transaction
	tx := s.DB.Begin()
	if tx.Error != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", tx.Error)
	}

	// Delete existing mock contents for this URL
	if err := tx.Where("url_id = ?", urlID).Delete(&models.MockContent{}).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to delete existing mock contents for url ID %d: %w", urlID, err)
	}

	// Create new mock contents
	// Ensure UrlID is set and ID is cleared for all items to be created
	for i := range mockContents {
		mockContents[i].UrlID = urlID
		mockContents[i].ID = 0 // Clear ID to ensure GORM creates new records
		mockContents[i].BaseModel.ID = 0
	}

	if len(mockContents) > 0 {
		if err := tx.Create(&mockContents).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create new mock contents for url ID %d: %w", urlID, err)
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return mockContents, nil
}

// SimulateLatency introduces a delay.
func (s *MockContentService) SimulateLatency(latencyMillis int64) {
	if latencyMillis > 0 {
		time.Sleep(time.Duration(latencyMillis) * time.Millisecond)
	}
}

// SelectRandomMockContent selects a mock content from a list based on weighted randomness.
// MockContents with higher 'Randomness' value have a higher chance of being selected.
func (s *MockContentService) SelectRandomMockContent(mockContents []models.MockContent) *models.MockContent {
	if len(mockContents) == 0 {
		return nil
	}

	// If all have zero randomness, or only one item, return the first one or a random one equally.
	// For simplicity, if no randomness is specified or applicable, return the first.
	// A more robust approach would be to check if all randomness values are 0.
	totalRandomness := int64(0)
	// First pass: sanitize negative randomness and calculate totalRandomness
	for i := range mockContents {
		if mockContents[i].Randomness < 0 {
			mockContents[i].Randomness = 0 // Modify the actual element in the slice
		}
		totalRandomness += mockContents[i].Randomness
	}

	// If total randomness is 0, all items are equally likely (or no items with randomness > 0).
	// In this case, pick one at random (uniformly).
	if totalRandomness == 0 {
		if len(mockContents) > 0 {
			// Seed rand if not already done globally, for better randomness.
			// rand.Seed(time.Now().UnixNano()) // Deprecated in Go 1.20+
			// Go 1.20+ automatically seeds. For older versions, seeding might be needed in main.
			return &mockContents[rand.Intn(len(mockContents))]
		}
		return nil
	}

	// Weighted random selection
	r := rand.Int63n(totalRandomness) // Generates a number between 0 and totalRandomness-1
	currentSum := int64(0)
	for i := range mockContents { // Iterate by index to use potentially modified values
		currentSum += mockContents[i].Randomness
		if r < currentSum {
			return &mockContents[i]
		}
	}

	// Fallback, should not be reached if totalRandomness > 0 and list is not empty.
	// But as a safeguard, return a random element if any exist (handles edge case of all zero after sanitizing)
	if len(mockContents) > 0 {
		return &mockContents[rand.Intn(len(mockContents))] // Fallback to uniform random if loop fails
	}
	return nil
}

// GetMockContentsByUrlID retrieves all mock contents for a given URL ID.
func (s *MockContentService) GetMockContentsByUrlID(urlID uint) ([]models.MockContent, error) {
	var contents []models.MockContent
	if err := s.DB.Where("url_id = ?", urlID).Find(&contents).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve mock contents for url ID %d: %w", urlID, err)
	}
	return contents, nil
}

// GetMockContentByID retrieves a single mock content by its ID.
func (s *MockContentService) GetMockContentByID(id uint) (*models.MockContent, error) {
	var content models.MockContent
	if err := s.DB.First(&content, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("mock content with ID %d not found: %w", id, err)
		}
		return nil, fmt.Errorf("failed to retrieve mock content with ID %d: %w", id, err)
	}
	return &content, nil
}

// DeleteMockContent deletes a mock content by its ID.
func (s *MockContentService) DeleteMockContent(id uint) error {
	result := s.DB.Delete(&models.MockContent{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete mock content with ID %d: %w", id, result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("mock content with ID %d not found for deletion: %w", id, gorm.ErrRecordNotFound)
	}
	return nil
}

// CreateMockContent creates a single mock content.
func (s *MockContentService) CreateMockContent(content *models.MockContent) error {
    if content == nil {
        return fmt.Errorf("mock content data cannot be nil")
    }
    // Ensure ID is zero for GORM to treat as new record
    content.ID = 0
    content.BaseModel.ID = 0

    if err := s.DB.Create(content).Error; err != nil {
        return fmt.Errorf("failed to create mock content: %w", err)
    }
    return nil
}

// UpdateMockContent updates a single mock content.
func (s *MockContentService) UpdateMockContent(content *models.MockContent) error {
    if content == nil || content.ID == 0 {
        return fmt.Errorf("mock content data is invalid or ID is missing for update")
    }
    if err := s.DB.Save(content).Error; err != nil {
        return fmt.Errorf("failed to update mock content with ID %d: %w", content.ID, err)
    }
    return nil
}
