package models

import (
	"database/sql"
	"time"
)

// RequestLog stores information about incoming requests
type RequestLog struct {
	ID        uint          `gorm:"primaryKey" json:"id"` // GORM's default ID
	IPAddress string        `json:"ip_address"`
	Timestamp time.Time     `json:"timestamp"`  // Timestamp of the request
	UrlID     sql.NullInt64 `json:"url_id"`     // Foreign key to Url, nullable if request doesn't match a defined Url
	ProjectID uint          `json:"project_id"` // To associate log with a project, even if UrlID is null
	Method    string        `json:"method"`     // HTTP method (GET, POST, etc.)
	Status    int           `json:"status"`     // HTTP status code returned
	URL       string        `json:"url"`        // The full requested URL
	IsProxied bool          `json:"is_proxied"` // True if the request was handled by the forward proxy
	CreatedAt time.Time     `json:"created_at"` // GORM will automatically manage this like @CreatedDate
}
