package models

// ForwardProxy represents the forward proxy settings for a project
type ForwardProxy struct {
	BaseModel
	Domain    string  `json:"domain"`
	ProjectID uint    `gorm:"unique;not null" json:"project_id"` // Foreign key for Project
	Project   Project `json:"project,omitempty"`                 // Belongs to Project
}
