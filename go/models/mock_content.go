package models

// MockContent represents the actual mock response data for a URL
type MockContent struct {
	BaseModel
	Randomness  int64  `gorm:"default:0;not null"` // For weighted random responses
	Latency     int64  `gorm:"default:0;not null"` // Latency in milliseconds
	Description string
	Name        string `gorm:"not null"`
	Data        string `gorm:"type:text;not null"` // The actual mock response body (e.g., JSON, XML)
	UrlID       uint   `gorm:"not null"`           // Foreign key for Url
	URL         Url    // Belongs to Url
}
