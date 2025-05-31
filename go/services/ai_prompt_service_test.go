package services_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/genai"

	"mockapi/config"
	"mockapi/services"
)

// MockModelsService implements the ModelsServiceInterface for testing.
type MockModelsService struct {
	GenerateContentFunc func(ctx context.Context, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error)
}

func (m *MockModelsService) GenerateContent(ctx context.Context, model string, contents []*genai.Content, config *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
	if m.GenerateContentFunc != nil {
		return m.GenerateContentFunc(ctx, model, contents, config)
	}
	return nil, fmt.Errorf("MockModelsService.GenerateContentFunc not implemented")
}

func TestAIPromptService_GetAIResponse(t *testing.T) {
	ctx := context.Background()
	dummyConfig := config.Config{GeminiModelName: "test-model"}

	t.Run("success_case", func(t *testing.T) {
		mockModelsSvc := &MockModelsService{
			GenerateContentFunc: func(c context.Context, modelName string, contents []*genai.Content, cfg *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
				// Basic checks on inputs
				require.Len(t, contents, 1, "Should receive one content item")
				require.Len(t, contents[0].Parts, 1, "Content item should have one part")

				inputPart := contents[0].Parts[0] // This is *genai.Part
				require.NotNil(t, inputPart, "Input part should not be nil")
				// genai.Part is a struct, access .Text field directly.
				assert.Equal(t, "test prompt", inputPart.Text)
				assert.Equal(t, "test-model", modelName)

				// Mock a fuller GenerateContentResponse structure
				respTextPart := genai.NewPartFromText("Hello from mock AI")
				return &genai.GenerateContentResponse{
					Candidates: []*genai.Candidate{
						{
							Content: &genai.Content{
								Parts: []*genai.Part{respTextPart},
								Role:  "model",
							},
							FinishReason: genai.FinishReasonStop,
							TokenCount:   12, // Example field for testing full marshalling
						},
					},
					// PromptFeedback can also be part of the response, add if needed for tests
				}, nil
			},
		}

		service, err := services.NewAIPromptService(dummyConfig, mockModelsSvc)
		require.NoError(t, err)
		require.NotNil(t, service)

		response, err := service.GetAIResponse(ctx, "test prompt")
		require.NoError(t, err)
		assert.NotNil(t, response)

		// Assertions for the fully marshalled GenerateContentResponse structure
		candidates, ok := response["candidates"].([]interface{})
		require.True(t, ok, "response should have 'candidates' field as a slice")
		require.Len(t, candidates, 1, "should be one candidate")

		candidateMap, ok := candidates[0].(map[string]interface{})
		require.True(t, ok, "candidate should be a map")

		assert.Equal(t, string(genai.FinishReasonStop), candidateMap["finishReason"])
		// JSON unmarshalling typically converts numbers to float64
		assert.Equal(t, float64(12), candidateMap["tokenCount"])

		contentMap, ok := candidateMap["content"].(map[string]interface{})
		require.True(t, ok, "candidate should have 'content' field as a map")

		partsList, ok := contentMap["parts"].([]interface{})
		require.True(t, ok, "content should have 'parts' field as a slice")
		require.Len(t, partsList, 1, "parts slice should have one element")

		partMap, ok := partsList[0].(map[string]interface{})
		require.True(t, ok, "part element should be a map")
		assert.Equal(t, "Hello from mock AI", partMap["text"])
		// Role is part of Content, not Part struct directly in the same way for simple text part.
		// The 'Role' field is on the 'Content' struct itself.
		assert.Equal(t, "model", contentMap["role"])

	})

	t.Run("error_case_service_fails", func(t *testing.T) {
		mockModelsSvc := &MockModelsService{
			GenerateContentFunc: func(c context.Context, modelName string, contents []*genai.Content, cfg *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
				return nil, errors.New("mock Gemini API error")
			},
		}
		service, err := services.NewAIPromptService(dummyConfig, mockModelsSvc)
		require.NoError(t, err)

		_, err = service.GetAIResponse(ctx, "test prompt")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "gemini API call failed")
		assert.Contains(t, err.Error(), "mock Gemini API error")
	})

	t.Run("error_case_empty_candidates_in_response", func(t *testing.T) {
		mockModelsSvc := &MockModelsService{
			GenerateContentFunc: func(c context.Context, modelName string, contents []*genai.Content, cfg *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
				return &genai.GenerateContentResponse{Candidates: []*genai.Candidate{}}, nil
			},
		}
		service, err := services.NewAIPromptService(dummyConfig, mockModelsSvc)
		require.NoError(t, err)

		response, err := service.GetAIResponse(ctx, "test prompt")
		require.NoError(t, err) // Marshalling an empty response should not error
		assert.NotNil(t, response)

		// Check that candidates field is present and empty or nil
		candidates, ok := response["candidates"]
		if ok && candidates != nil { // It might be present as an empty list [] or nil if omitempty was used by genai's MarshalJSON
			candidateList, isList := candidates.([]interface{})
			assert.True(t, isList, "candidates field should be a list if present and not nil")
			assert.Len(t, candidateList, 0, "candidates list should be empty")
		} else {
			// If candidates is nil or not present, that's also acceptable for an empty response
			assert.True(t, candidates == nil || !ok, "candidates field should be nil or not present")
		}
	})

	t.Run("error_case_nil_response_from_service_call", func(t *testing.T) {
		mockModelsSvc := &MockModelsService{
			GenerateContentFunc: func(c context.Context, modelName string, contents []*genai.Content, cfg *genai.GenerateContentConfig) (*genai.GenerateContentResponse, error) {
				return nil, nil // Nil response, nil error
			},
		}
		service, err := services.NewAIPromptService(dummyConfig, mockModelsSvc)
		require.NoError(t, err)

		_, err = service.GetAIResponse(ctx, "test prompt")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "received nil response from Gemini API")
	})

	t.Run("constructor_error_nil_service", func(t *testing.T) {
		_, err := services.NewAIPromptService(dummyConfig, nil)
		require.Error(t, err)
		assert.EqualError(t, err, "ModelsServiceInterface cannot be nil for AIPromptService")
	})
}
