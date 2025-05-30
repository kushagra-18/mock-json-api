package services

import (
	"context"
	"encoding/json" // Will be used for marshalling the full response
	"fmt"
	"log"
	// "strings" // No longer needed after switching to full response marshalling

	// "google.golang.org/api/option" // Not used here if client is passed in
	"google.golang.org/genai"

	"mockapi/config"
)

// AIPromptService handles interactions with the Gemini AI model.
type AIPromptService struct {
	modelsService ModelsServiceInterface // Use the new interface for client.Models
	modelName     string               // Store the model name
}

// NewAIPromptService creates a new AIPromptService.
// It now accepts ModelsServiceInterface for better testability.
func NewAIPromptService(cfg config.Config, modelsSvc ModelsServiceInterface) (*AIPromptService, error) {
	if modelsSvc == nil {
		return nil, fmt.Errorf("ModelsServiceInterface cannot be nil for AIPromptService")
	}

	modelName := cfg.GeminiModelName
	if modelName == "" {
		log.Println("Warning: Gemini model name not configured in cfg, using default 'gemini-1.5-flash-latest'")
		modelName = "gemini-1.5-flash-latest"
	}

	return &AIPromptService{modelsService: modelsSvc, modelName: modelName}, nil
}

// GetAIResponse sends a prompt to the Gemini model and returns the response as a map.
func (s *AIPromptService) GetAIResponse(ctx context.Context, prompt string) (map[string]interface{}, error) {
	if s.modelsService == nil {
		return nil, fmt.Errorf("models service is not initialized in AIPromptService")
	}

	// Construct content for client.Models.GenerateContent
	part := genai.NewPartFromText(prompt) // Returns *Part
	// NewContentFromParts expects []*Part, and genai.RoleUser is "user"
	content := genai.NewContentFromParts([]*genai.Part{part}, genai.RoleUser)
	contents := []*genai.Content{content}

	// Call GenerateContent on the ModelsServiceInterface
	resp, err := s.modelsService.GenerateContent(ctx, s.modelName, contents, nil) // Pass nil for GenerateContentConfig for now
	if err != nil {
		log.Printf("Error from Gemini API call: %v", err)
		return nil, fmt.Errorf("gemini API call failed: %w", err)
	}

	if resp == nil {
		return nil, fmt.Errorf("received nil response from Gemini API")
	}

	var jsonResponse map[string]interface{}
	responseBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling Gemini response: %v", err)
		return nil, fmt.Errorf("failed to marshal gemini response: %w", err)
	}
	if err := json.Unmarshal(responseBytes, &jsonResponse); err != nil {
		log.Printf("Error unmarshalling Gemini response bytes to map: %v", err)
		return nil, fmt.Errorf("failed to unmarshal gemini response to map: %w", err)
	}
	return jsonResponse, nil
}
