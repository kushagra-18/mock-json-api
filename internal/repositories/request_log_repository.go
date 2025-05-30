package repositories

import (
	"errors"
	"go-gin-gorm-api/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// RequestLogRepository defines the interface for request log data operations.
type RequestLogRepository interface {
	CreateRequestLog(requestLog *models.RequestLog) error
	GetRequestLogByID(id uint) (*models.RequestLog, error)
	GetRequestLogsByProjectID(projectID uint) ([]models.RequestLog, error)
	GetAllRequestLogs() ([]models.RequestLog, error)
	// Add other query methods here if needed
}

// requestLogRepository implements RequestLogRepository with GORM.
type requestLogRepository struct {
	db *gorm.DB
}

// NewRequestLogRepository creates a new instance of requestLogRepository.
func NewRequestLogRepository(db *gorm.DB) RequestLogRepository {
	return &requestLogRepository{db: db}
}

// CreateRequestLog creates a new request log in the database.
func (r *requestLogRepository) CreateRequestLog(requestLog *models.RequestLog) error {
	return r.db.Create(requestLog).Error
}

// GetRequestLogByID retrieves a request log by its ID.
// Preloads Project association.
func (r *requestLogRepository) GetRequestLogByID(id uint) (*models.RequestLog, error) {
	var requestLog models.RequestLog
	err := r.db.Preload(clause.Associations).First(&requestLog, id).Error // "Project" or models.RequestLog{}.Project
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or return a specific "not found" error
		}
		return nil, err
	}
	return &requestLog, nil
}

// GetRequestLogsByProjectID retrieves all request logs for a given projectID.
// Preloads Project association.
func (r *requestLogRepository) GetRequestLogsByProjectID(projectID uint) ([]models.RequestLog, error) {
	var requestLogs []models.RequestLog
	err := r.db.Preload(clause.Associations).Where("project_id = ?", projectID).Find(&requestLogs).Error
	return requestLogs, err
}

// GetAllRequestLogs retrieves all request logs from the database.
// Preloads Project association.
func (r *requestLogRepository) GetAllRequestLogs() ([]models.RequestLog, error) {
	var requestLogs []models.RequestLog
	err := r.db.Preload(clause.Associations).Find(&requestLogs).Error
	return requestLogs, err
}
