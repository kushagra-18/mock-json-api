package models

// Team represents a team in the system
type Team struct {
	BaseModel
	Name     string    `gorm:"unique;not null"`
	Slug     string    `gorm:"unique;not null"`
	Projects []Project `gorm:"foreignKey:TeamID"`
}
