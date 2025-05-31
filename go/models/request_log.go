package models

import (
	"database/sql"
	"time"
)

// RequestLog stores information about incoming requests
type RequestLog struct {
	ID        uint `gorm:"primaryKey"` // GORM's default ID
	IPAddress string
	Timestamp time.Time     // Timestamp of the request
	UrlID     sql.NullInt64 // Foreign key to Url, nullable if request doesn't match a defined Url
	ProjectID uint          // To associate log with a project, even if UrlID is null
	Method    string        // HTTP method (GET, POST, etc.)
	Status    int           // HTTP status code returned
	URL       string        // The full requested URL
	IsProxied bool          // True if the request was handled by the forward proxy
	CreatedAt time.Time     // GORM will automatically manage this like @CreatedDate
}
