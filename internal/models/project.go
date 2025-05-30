package models

import (
	"time"

	"gorm.io/gorm"
)

// Project represents a project in the system.
type Project struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Name      string         `gorm:"uniqueIndex:idx_project_name_team_id;not null" json:"name"`
	Slug      string         `gorm:"uniqueIndex:idx_project_slug_team_id;not null" json:"slug"`
	TeamID    uint           `gorm:"uniqueIndex:idx_project_name_team_id;uniqueIndex:idx_project_slug_team_id;not null" json:"team_id"`
	Team      *Team          `gorm:"foreignKey:TeamID" json:"team,omitempty"`
}

// TableName specifies the table name for the Project model.
func (Project) TableName() string {
	return "projects"
}
