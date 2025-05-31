package models

import (
	"database/sql"
)

// Url represents a URL endpoint within a project that can be mocked
type Url struct {
	BaseModel
	Description  string        `json:"description"`
	Name         string        `gorm:"not null" json:"name"`
	Requests     sql.NullInt64 `json:"requests"`                                         // Nullable integer for number of requests
	Time         sql.NullInt64 `json:"time"`                                             // Nullable integer for response time (e.g., in ms)
	URL          string        `gorm:"not null;index:idx_url_project,unique" json:"url"` // URL path, unique per project
	Status       StatusCode    `gorm:"type:varchar(50);not null" json:"status"`
	ProjectID    uint          `gorm:"index:idx_url_project,unique" json:"project_id"`  // Foreign key for Project, part of composite unique index
	Project      Project       `json:"project,omitempty"`                               // Belongs to Project
	MockContents []MockContent `gorm:"foreignKey:UrlID" json:"mock_contents,omitempty"` // Has many MockContents
}
