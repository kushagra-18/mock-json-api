package models

import (
	"time"

	"gorm.io/gorm"
)

// Team represents a team in the system.
type Team struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"` // Use json:"-" to exclude from JSON output
	Name      string         `gorm:"uniqueIndex;not null" json:"name"`
	Slug      string         `gorm:"uniqueIndex;not null" json:"slug"`
	Projects  []Project      `gorm:"foreignKey:TeamID" json:"projects,omitempty"`
}

// TableName specifies the table name for the Team model.
func (Team) TableName() string {
	return "teams"
}
