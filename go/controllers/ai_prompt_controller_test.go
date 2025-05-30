package controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"mockapi/controllers"
	"mockapi/dtos"
	"mockapi/services" // For the interface
	// "mockapi/utils" // Not directly used by test functions
)

var _ services.AIPromptServiceInterface // Force services import usage

// MockAIPromptService implements the AIPromptServiceInterface for testing.
type MockAIPromptService struct {
	GetAIResponseFunc func(ctx context.Context, prompt string) (map[string]interface{}, error)
	// CloseFunc is not needed as the interface no longer has Close()
}

func (m *MockAIPromptService) GetAIResponse(ctx context.Context, prompt string) (map[string]interface{}, error) {
	if m.GetAIResponseFunc != nil {
		return m.GetAIResponseFunc(ctx, prompt)
	}
	return nil, errors.New("MockAIPromptService.GetAIResponseFunc not implemented")
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New() // Use gin.New() instead of gin.Default() for tests to avoid default middleware
	return router
}

func TestAIPromptController_HandleAIPrompt(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		mockService := &MockAIPromptService{
			GetAIResponseFunc: func(ctx context.Context, prompt string) (map[string]interface{}, error) {
				assert.Equal(t, "test prompt", prompt)
				return map[string]interface{}{"text": "AI response to " + prompt}, nil
			},
		}
		controller := controllers.NewAIPromptController(mockService)

		router := setupTestRouter()
		router.POST("/ai/prompt", controller.HandleAIPrompt)

		reqBody := dtos.AIPromptRequestDTO{Prompt: "test prompt"}
		jsonReqBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, "/ai/prompt", bytes.NewBuffer(jsonReqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)

		var respBody map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &respBody)
		require.NoError(t, err)

		assert.Equal(t, "success", respBody["status"])
		data, ok := respBody["data"].(map[string]interface{})
		require.True(t, ok)
		assert.Equal(t, "AI response to test prompt", data["text"])
	})

	t.Run("invalid_request_bad_json", func(t *testing.T) {
		// Provide a valid, non-nil mockService, even if its methods won't be called.
		mockService := &MockAIPromptService{}
		controller := controllers.NewAIPromptController(mockService)

		router := setupTestRouter()
		router.POST("/ai/prompt", controller.HandleAIPrompt)

		req, _ := http.NewRequest(http.MethodPost, "/ai/prompt", bytes.NewBufferString("{bad json"))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		// utils.ErrorResponse format
		var errResp map[string]interface{}
		err := json.Unmarshal(rr.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "error", errResp["status"])
		assert.Contains(t, errResp["message"], "Invalid request payload")
	})

	t.Run("invalid_request_empty_prompt", func(t *testing.T) {
		// Provide a valid, non-nil mockService.
		mockService := &MockAIPromptService{}
		controller := controllers.NewAIPromptController(mockService)

		router := setupTestRouter()
		router.POST("/ai/prompt", controller.HandleAIPrompt)

		reqBody := dtos.AIPromptRequestDTO{Prompt: ""} // Empty prompt
		jsonReqBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, "/ai/prompt", bytes.NewBuffer(jsonReqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		var errResp map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "error", errResp["status"])
		assert.Contains(t, errResp["message"], "Invalid request payload") // Gin binding error due to "min=1"
	})


	t.Run("service_error_case", func(t *testing.T) {
		mockService := &MockAIPromptService{
			GetAIResponseFunc: func(ctx context.Context, prompt string) (map[string]interface{}, error) {
				return nil, errors.New("mock service error")
			},
		}
		controller := controllers.NewAIPromptController(mockService)

		router := setupTestRouter()
		router.POST("/ai/prompt", controller.HandleAIPrompt)

		reqBody := dtos.AIPromptRequestDTO{Prompt: "test prompt"}
		jsonReqBody, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, _ := http.NewRequest(http.MethodPost, "/ai/prompt", bytes.NewBuffer(jsonReqBody))
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		var errResp map[string]interface{}
		err = json.Unmarshal(rr.Body.Bytes(), &errResp)
		require.NoError(t, err)
		assert.Equal(t, "error", errResp["status"])
		assert.Contains(t, errResp["message"], "Failed to get AI response: mock service error")
	})

	t.Run("nil_service_in_controller_constructor_panic", func(t *testing.T) {
		// The constructor NewAIPromptController now panics if service is nil.
		assert.PanicsWithValue(t, "AIPromptServiceInterface cannot be nil in NewAIPromptController", func() {
			controllers.NewAIPromptController(nil)
		}, "Constructor should panic if service is nil")
	})
}
