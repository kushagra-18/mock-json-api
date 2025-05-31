package dtos

// AIPromptRequestDTO defines the structure for AI prompt requests.
type AIPromptRequestDTO struct {
	Prompt string `json:"prompt" binding:"required,min=1"` // Added min=1 to ensure not just empty string after trim
}
