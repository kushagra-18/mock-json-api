package models

import (
	"time"
)

// RequestLog represents a request log in the system.
type RequestLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	URL       string    `gorm:"column:url;not null" json:"url"` // Matched json tag with existing annotation
	Method    string    `gorm:"not null" json:"method"`
	IP        string    `gorm:"column:ip_address" json:"ip_address"` // Matched json tag with existing annotation
	Status    int       `gorm:"not null" json:"status"`
	ProjectID uint      `gorm:"not null" json:"project_id"`
	Project   *Project  `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
}

// TableName specifies the table name for the RequestLog model.
func (RequestLog) TableName() string {
	return "request_logs"
}
