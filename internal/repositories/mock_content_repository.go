package repositories

import (
	"errors"
	"go-gin-gorm-api/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MockContentRepository defines the interface for mock content data operations.
type MockContentRepository interface {
	CreateMockContent(mockContent *models.MockContent) error
	GetMockContentByID(id uint) (*models.MockContent, error)
	GetMockContentsByUrlID(urlID uint) ([]models.MockContent, error)
	UpdateMockContent(mockContent *models.MockContent) error
	DeleteMockContent(id uint) error
}

// mockContentRepository implements MockContentRepository with GORM.
type mockContentRepository struct {
	db *gorm.DB
}

// NewMockContentRepository creates a new instance of mockContentRepository.
func NewMockContentRepository(db *gorm.DB) MockContentRepository {
	return &mockContentRepository{db: db}
}

// CreateMockContent creates new mock content in the database.
func (r *mockContentRepository) CreateMockContent(mockContent *models.MockContent) error {
	return r.db.Create(mockContent).Error
}

// GetMockContentByID retrieves mock content by its ID, preloading Url.
func (r *mockContentRepository) GetMockContentByID(id uint) (*models.MockContent, error) {
	var mockContent models.MockContent
	err := r.db.Preload(clause.Associations).First(&mockContent, id).Error // "Url" or models.Url{}.Url
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or return a specific "not found" error
		}
		return nil, err
	}
	return &mockContent, nil
}

// GetMockContentsByUrlID retrieves all mock content for a given urlID, preloading Url.
func (r *mockContentRepository) GetMockContentsByUrlID(urlID uint) ([]models.MockContent, error) {
	var mockContents []models.MockContent
	// Preload("Url") or Preload(clause.Associations) should work if gorm tags are correct
	err := r.db.Preload(clause.Associations).Where("url_id = ?", urlID).Find(&mockContents).Error
	return mockContents, err
}

// UpdateMockContent updates existing mock content in the database.
func (r *mockContentRepository) UpdateMockContent(mockContent *models.MockContent) error {
	return r.db.Save(mockContent).Error
}

// DeleteMockContent soft deletes mock content by its ID.
func (r *mockContentRepository) DeleteMockContent(id uint) error {
	return r.db.Delete(&models.MockContent{}, id).Error
}
