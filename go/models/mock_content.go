package models

// MockContent represents the actual mock response data for a URL
type MockContent struct {
	BaseModel
	Randomness  int64  `gorm:"default:0;not null" json:"randomness"` // For weighted random responses
	Latency     int64  `gorm:"default:0;not null" json:"latency"`    // Latency in milliseconds
	Description string `json:"description"`
	Name        string `gorm:"not null" json:"name"`
	Data        string `gorm:"type:text;not null" json:"data"` // The actual mock response body (e.g., JSON, XML)
	UrlID       uint   `gorm:"not null" json:"url_id"`         // Foreign key for Url
	URL         Url    `json:"url,omitempty"`                  // Belongs to Url
}
