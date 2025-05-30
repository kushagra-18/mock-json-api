package models

// ForwardProxy represents the forward proxy settings for a project
type ForwardProxy struct {
	BaseModel
	Domain    string
	ProjectID uint `gorm:"unique;not null"` // Foreign key for Project
	Project   Project // Belongs to Project
}
