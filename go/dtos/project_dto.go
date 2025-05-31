package dtos

// CreateProjectDTO is used for creating a new project.
// Pointers are used for optional fields. Validation tags ensure required fields are present.
type CreateProjectDTO struct {
	Name        *string `json:"name"`                                 // Optional, can be derived from slug or set to a default
	Slug        *string `json:"slug" binding:"required,min=3,max=50"` // Slug is typically required and has length constraints
	Description *string `json:"description"`
	TeamID      *uint   `json:"team_id"` // Optional: If not provided, might use a default team or user's primary team
}
