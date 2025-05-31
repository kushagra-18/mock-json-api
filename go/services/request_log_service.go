package services

import (
	"fmt"
	"log" // For placeholder Pusher event

	"gorm.io/gorm"
	"mockapi/models" // Assuming module name is mockapi
)

// RequestLogService handles business logic related to request logs.
// PusherService integration is stubbed out for now.
type RequestLogService struct {
	DB *gorm.DB
	// PusherService *PusherService // Uncomment and use when PusherService is implemented
}

// NewRequestLogService creates a new RequestLogService.
// func NewRequestLogService(db *gorm.DB, pusherService *PusherService) *RequestLogService {
// 	return &RequestLogService{DB: db, PusherService: pusherService}
// }
func NewRequestLogService(db *gorm.DB) *RequestLogService {
	return &RequestLogService{DB: db}
}

// SaveRequestLog saves a new request log.
func (s *RequestLogService) SaveRequestLog(requestLog *models.RequestLog) error {
	if requestLog == nil {
		return fmt.Errorf("request log data cannot be nil")
	}

	// The RequestLog model has its own ID (uint `gorm:"primaryKey"`) and CreatedAt,
	// so GORM should handle these automatically.

	if err := s.DB.Create(requestLog).Error; err != nil {
		return fmt.Errorf("failed to save request log: %w", err)
	}

	// After saving, attempt to emit a Pusher event.
	// This is a placeholder for the actual Pusher integration.
	if requestLog.ProjectID != 0 { // Only emit if associated with a project
		s.EmitPusherEvent(requestLog.ProjectID, "new_request", requestLog)
	}

	return nil
}

// GetLogsByProjectID retrieves request logs for a specific project with pagination.
func (s *RequestLogService) GetLogsByProjectID(projectID uint, limit, offset int) ([]models.RequestLog, error) {
	var logs []models.RequestLog
	query := s.DB.Where("project_id = ?", projectID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&logs).Error; err != nil {
		return nil, fmt.Errorf("failed to retrieve logs for project ID %d: %w", projectID, err)
	}
	return logs, nil
}

// EmitPusherEvent is a placeholder for sending data via Pusher/WebSockets.
// In a real implementation, this would interact with a Pusher client or similar.
func (s *RequestLogService) EmitPusherEvent(projectID uint, eventName string, data interface{}) {
	// Placeholder: Log the event instead of sending it via Pusher.
	// In a real app, this would use the PusherService.
	// Example: s.PusherService.Trigger(fmt.Sprintf("project-%d", projectID), eventName, data)
	log.Printf("Pusher Event Emitted (Placeholder): ProjectID=%d, Event=%s, Data=%v\n", projectID, eventName, data)
}

// GetLogByID retrieves a single request log by its ID.
func (s *RequestLogService) GetLogByID(id uint) (*models.RequestLog, error) {
    var requestLog models.RequestLog
    if err := s.DB.First(&requestLog, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, fmt.Errorf("request log with ID %d not found: %w", id, err)
        }
        return nil, fmt.Errorf("failed to retrieve request log with ID %d: %w", id, err)
    }
    return &requestLog, nil
}

// DeleteLogsByProjectID deletes all request logs for a given project ID.
// This can be a potentially long-running operation and might need care in a production system.
func (s *RequestLogService) DeleteLogsByProjectID(projectID uint) (int64, error) {
    result := s.DB.Where("project_id = ?", projectID).Delete(&models.RequestLog{})
    if result.Error != nil {
        return 0, fmt.Errorf("failed to delete request logs for project ID %d: %w", projectID, result.Error)
    }
    return result.RowsAffected, nil
}
