package models

import (
	"time"

	"gorm.io/gorm"
)

// Url represents a URL in the system.
type Url struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	Description     string         `json:"description"`
	Name            string         `gorm:"uniqueIndex:idx_url_name_project_id;not null" json:"name"`
	URL             string         `gorm:"column:url;uniqueIndex:idx_url_url_project_id;not null" json:"url"`
	ProjectID       uint           `gorm:"uniqueIndex:idx_url_name_project_id;uniqueIndex:idx_url_url_project_id;not null" json:"project_id"`
	Project         *Project       `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	MockContentList []MockContent  `gorm:"foreignKey:URLID" json:"mock_content_list,omitempty"`
}

// TableName specifies the table name for the Url model.
func (Url) TableName() string {
	return "urls"
}
