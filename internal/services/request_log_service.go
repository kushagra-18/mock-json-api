package services

import (
	"go-gin-gorm-api/internal/models"
	"go-gin-gorm-api/internal/repositories"
	"log"
	"time"
)

// RequestLogService defines the interface for request log-related business logic.
type RequestLogService interface {
	CreateRequestLog(requestLog *models.RequestLog) error
	SaveRequestLogAsync(url, method, ipAddress string, status int, projectID uint)
	EmitPusherEventAsync(method, url string, projectID uint, status int) // Placeholder
}

// requestLogService implements RequestLogService.
type requestLogService struct {
	logRepo     repositories.RequestLogRepository
	projectRepo repositories.ProjectRepository // For potential validation or project-related logic
	// Pusher client/service would be added here later
}

// NewRequestLogService creates a new instance of RequestLogService.
func NewRequestLogService(logRepo repositories.RequestLogRepository, projectRepo repositories.ProjectRepository) RequestLogService {
	return &requestLogService{logRepo: logRepo, projectRepo: projectRepo}
}

// CreateRequestLog creates a new request log synchronously.
func (s *requestLogService) CreateRequestLog(requestLog *models.RequestLog) error {
	// Consider validating projectID if necessary:
	// project, err := s.projectRepo.GetProjectByID(requestLog.ProjectID)
	// if err != nil {
	// 	return fmt.Errorf("error validating project ID for request log: %w", err)
	// }
	// if project == nil {
	// 	return fmt.Errorf("project with ID %d not found for request log", requestLog.ProjectID)
	// }
	// RequestLog model's CreatedAt is automatically handled by GORM if it's time.Time
	return s.logRepo.CreateRequestLog(requestLog)
}

// SaveRequestLogAsync saves a request log asynchronously.
func (s *requestLogService) SaveRequestLogAsync(url, method, ipAddress string, status int, projectID uint) {
	go func() {
		// Optional: Validate projectID exists before creating the log.
		// This might be less critical for an async log to not block the main request flow,
		// but depends on requirements. If validation fails, the log simply isn't saved.
		// For example:
		// _, err := s.projectRepo.GetProjectByID(projectID)
		// if err != nil {
		// 	log.Printf("Error validating project ID %d for async request log: %v. Log not saved.", projectID, err)
		// 	return
		// }

		requestLog := &models.RequestLog{
			URL:       url,
			Method:    method,
			IP:        ipAddress,
			Status:    status,
			ProjectID: projectID,
			CreatedAt: time.Now(), // Explicitly set for async or rely on GORM if configured
		}
		if err := s.logRepo.CreateRequestLog(requestLog); err != nil {
			log.Printf("Error saving request log asynchronously: %v", err)
		}
	}()
}

// EmitPusherEventAsync is a placeholder for emitting events via Pusher.
// It currently logs the event details and runs in a goroutine.
func (s *requestLogService) EmitPusherEventAsync(method, url string, projectID uint, status int) {
	go func() {
		// Actual Pusher integration would go here.
		log.Printf("Pusher event triggered (placeholder): method=%s, url=%s, projectID=%d, status=%d", method, url, projectID, status)
		// Example:
		// client := getPusherClient() // Assuming a way to get an initialized Pusher client
		// eventData := map[string]string{
		// 	"method": method,
		// 	"url":    url,
		// 	"status": fmt.Sprintf("%d", status),
		// }
		// channelName := fmt.Sprintf("project-%d-logs", projectID)
		// err := client.Trigger(channelName, "new-log", eventData)
		// if err != nil {
		// 	log.Printf("Error emitting Pusher event: %v", err)
		// }
	}()
}
