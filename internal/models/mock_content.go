package models

import (
	"time"

	"gorm.io/gorm"
)

// MockContent represents mock content for a URL in the system.
type MockContent struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	URLID       uint           `gorm:"not null" json:"url_id"`
	Url         *Url           `gorm:"foreignKey:URLID" json:"url,omitempty"` // Changed from URL to Url to match struct name
	Randomness  int            `json:"randomness"`
	Latency     int64          `json:"latency"`
	Description string         `json:"description"`
	Name        string         `gorm:"not null" json:"name"`
	Data        string         `gorm:"type:text" json:"data"`
}

// TableName specifies the table name for the MockContent model.
func (MockContent) TableName() string {
	return "mock_content"
}
