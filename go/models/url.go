package models

import (
	"database/sql"
)

// Url represents a URL endpoint within a project that can be mocked
type Url struct {
	BaseModel
	Description  string
	Name         string        `gorm:"not null"`
	Requests     sql.NullInt64 // Nullable integer for number of requests
	Time         sql.NullInt64 // Nullable integer for response time (e.g., in ms)
	URL          string        `gorm:"not null;index:idx_url_project,unique"` // URL path, unique per project
	Status       StatusCode    `gorm:"type:varchar(50);not null"`
	ProjectID    uint          `gorm:"index:idx_url_project,unique"` // Foreign key for Project, part of composite unique index
	Project      Project       // Belongs to Project
	MockContents []MockContent `gorm:"foreignKey:UrlID"` // Has many MockContents
}
