package models

// Team represents a team in the system
type Team struct {
	BaseModel
	Name     string    `gorm:"unique;not null" json:"name"`
	Slug     string    `gorm:"unique;not null" json:"slug"`
	Projects []Project `gorm:"foreignKey:TeamID" json:"projects,omitempty"`
}
