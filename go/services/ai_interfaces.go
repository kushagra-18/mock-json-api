package services

import (
	"context"
	"google.golang.org/genai"
)

// GenerativeModelInterface defines the methods we use from *genai.GenerativeModel
// This allows for mocking the actual Gemini model in tests.
type GenerativeModelInterface interface {
	GenerateContent(ctx context.Context, parts ...genai.Part) (*genai.GenerateContentResponse, error)
	// Add other methods from *genai.GenerativeModel if used by AIPromptService in the future
}

// ModelsServiceInterface defines the methods we use from *genai.Client.Models
type ModelsServiceInterface interface {
	GenerateContent(ctx context.Context, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error)
	// Add other methods from *genai.Models if used
}


// AIPromptServiceInterface defines the methods for AIPromptService.
// This allows for mocking the service in controller tests.
type AIPromptServiceInterface interface {
	GetAIResponse(ctx context.Context, prompt string) (map[string]interface{}, error)
	// Close() method is removed as client lifecycle is managed externally.
}
