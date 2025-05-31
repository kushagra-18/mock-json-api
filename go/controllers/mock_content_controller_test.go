package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm" // For gorm.ErrRecordNotFound

	"mockapi/config"
	"mockapi/controllers"
	"mockapi/dtos"
	"mockapi/models"
	"mockapi/services" // Required for the mock service types
	// "mockapi/utils" // Not strictly needed for this version of the test file
)

// Helper struct to hold all mocks for controller tests
type controllerMocks struct {
	mockProjectSvc     *services.MockProjectService
	mockMcSvc          *services.MockMockContentService
	mockUrlSvc         *services.MockURLService
	mockReqLogSvc      *services.MockRequestLogService
	mockRedisSvc       *services.MockRedisService
	mockProxySvc       *services.MockProxyService
	mockFakerSvc       *services.MockFakerService
}

// setupTestRouterWithMocks initializes a Gin router and MockContentController with all mocks.
func setupTestRouterWithMocks(t *testing.T) (*gin.Engine, *controllerMocks, *controllers.MockContentController) {
	gin.SetMode(gin.TestMode) // Important for testing
	router := gin.Default()

	mocks := &controllerMocks{
		mockProjectSvc:     &services.MockProjectService{},
		mockMcSvc:          &services.MockMockContentService{},
		mockUrlSvc:         &services.MockURLService{},
		mockReqLogSvc:      &services.MockRequestLogService{},
		mockRedisSvc:       &services.MockRedisService{},
		mockProxySvc:       &services.MockProxyService{},
		mockFakerSvc:       &services.MockFakerService{},
	}

	// Basic config for tests
	cfg := config.Config{
		JWTSecretKey: "testsecret",
		// Fill other necessary config fields if they affect controller logic being tested
	}

	mcController := controllers.NewMockContentController(
		mocks.mockProjectSvc,
		mocks.mockMcSvc,
		mocks.mockUrlSvc,
		mocks.mockReqLogSvc,
		mocks.mockRedisSvc,
		mocks.mockProxySvc,
		mocks.mockFakerSvc,
		cfg,
	)

	return router, mocks, mcController
}

func TestMockContentController_SaveMockContent_WithDSL_Success(t *testing.T) {
	router, mocks, mcController := setupTestRouterWithMocks(t)

	// 1. Setup Mocks for dependent services (Project, URL, MockContent)
	mocks.mockProjectSvc.GetProjectBySlugFunc = func(slug string) (*models.Project, error) {
		if slug == "test-project" {
			return &models.Project{BaseModel: models.BaseModel{ID: 1}, Slug: "test-project"}, nil
		}
		return nil, gorm.ErrRecordNotFound
	}
	mocks.mockUrlSvc.FindByProjectIDAndURLFunc = func(projectID uint, urlPath string) (*models.Url, error) {
		return nil, gorm.ErrRecordNotFound // Assume URL doesn't exist yet
	}

	var capturedUrlArg *models.Url
	mocks.mockUrlSvc.CreateURLFunc = func(url *models.Url, projectID uint) error {
		capturedUrlArg = url
		url.ID = 123         // Assign an ID as the actual service would
		return nil
	}

	var capturedMockContents []models.MockContent
	mocks.mockMcSvc.SaveMockContentListFunc = func(mockContents []models.MockContent, urlID uint) ([]models.MockContent, error) {
		capturedMockContents = mockContents
		for i := range mockContents {
			mockContents[i].ID = uint(i + 1)
		}
		return mockContents, nil
	}

	// 2. Configure MockFakerService
	expectedDSL := "{{name.firstName}}"
	expectedProcessedData := "\"John Doe From Faker\""
	mocks.mockFakerSvc.ProcessDSLFunc = func(dsl string) (string, error) {
		if dsl == expectedDSL {
			return expectedProcessedData, nil
		}
		return "", fmt.Errorf("unexpected DSL: %s", dsl)
	}

	// 3. Setup Router and Endpoint
	// The controller instance mcController is already configured with mocks.
	// We need to register its methods with the router.
	router.POST("/mock/:projectSlug", mcController.SaveMockContent)

	// 4. Prepare Request
	// Corrected to use MockContentCreateDTO as per DTO definitions
	dslPayload := dtos.MockContentUrlDTO{
		URLData: struct {
			Description *string          `json:"description"`
			Name        string           `json:"name" binding:"required"`
			URL         string           `json:"url" binding:"required"`
			Status      models.StatusCode `json:"status" binding:"required"`
		}{
			Name:   "Test DSL URL",
			URL:    "/test-dsl-path",
			Status: models.StatusOK,
		},
		MockContentList: []dtos.MockContentCreateDTO{ // Corrected DTO type
			{
				Name:    "Test DSL Content Item",
				DslData: &expectedDSL,
				Data:    "", // Explicitly empty static data
			},
		},
	}
	bodyBytes, _ := json.Marshal(dslPayload)
	req, _ := http.NewRequest("POST", "/mock/test-project", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	// 5. Serve HTTP Request
	router.ServeHTTP(resp, req)

	// 6. Assertions
	if resp.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusCreated, resp.Code, resp.Body.String())
	}

	if capturedMockContents == nil {
		t.Fatalf("MockMockContentService.SaveMockContentListFunc was not called")
	}
	if len(capturedMockContents) != 1 {
		t.Fatalf("expected 1 mock content to be saved, got %d", len(capturedMockContents))
	}
	if capturedMockContents[0].Data != expectedProcessedData {
		t.Errorf("expected processed data '%s', got '%s'", expectedProcessedData, capturedMockContents[0].Data)
	}
	if capturedMockContents[0].Name != "Test DSL Content Item" {
		t.Errorf("unexpected mock content name: %s", capturedMockContents[0].Name)
	}

	if capturedUrlArg == nil {
		t.Fatalf("MockURLService.CreateURLFunc was not called")
	}
	if capturedUrlArg.URL != "/test-dsl-path" {
		t.Errorf("expected URL path '/test-dsl-path', got '%s'", capturedUrlArg.URL)
	}

	t.Log("SaveMockContent_WithDSL_Success completed. Response:", resp.Body.String())
}


func TestMockContentController_SaveMockContent_DSLError(t *testing.T) {
	router, mocks, mcController := setupTestRouterWithMocks(t)

	mocks.mockProjectSvc.GetProjectBySlugFunc = func(slug string) (*models.Project, error) {
		return &models.Project{BaseModel: models.BaseModel{ID: 1}, Slug: slug}, nil
	}
	mocks.mockUrlSvc.FindByProjectIDAndURLFunc = func(projectID uint, urlPath string) (*models.Url, error) {
		return nil, gorm.ErrRecordNotFound
	}
    // CreateURL and SaveMockContentList should not be called if DSL processing fails.

	expectedDSL := "{{name.firstName}}"
	expectedErrorMessage := "faker processing error from test"
	mocks.mockFakerSvc.ProcessDSLFunc = func(dsl string) (string, error) {
		if dsl == expectedDSL {
			return "", errors.New(expectedErrorMessage)
		}
		return "", fmt.Errorf("unexpected DSL in error test: %s", dsl)
	}

	router.POST("/mock/:projectSlug", mcController.SaveMockContent)

	// Corrected to use MockContentCreateDTO
	dslPayload := dtos.MockContentUrlDTO{
		URLData: struct {
			Description *string          `json:"description"`
			Name        string           `json:"name" binding:"required"`
			URL         string           `json:"url" binding:"required"`
			Status      models.StatusCode `json:"status" binding:"required"`
		}{
			Name:   "Test DSL URL Error",
			URL:    "/test-dsl-error",
			Status: models.StatusOK,
		},
		MockContentList: []dtos.MockContentCreateDTO{ // Corrected DTO type
			{
				Name:    "Test DSL Error Item",
				DslData: &expectedDSL,
			},
		},
	}
	bodyBytes, _ := json.Marshal(dslPayload)
	req, _ := http.NewRequest("POST", "/mock/test-project-error", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	if resp.Code != http.StatusInternalServerError {
		t.Errorf("expected status %d, got %d. Response: %s", http.StatusInternalServerError, resp.Code, resp.Body.String())
	}

	var errorResponse map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &errorResponse); err != nil {
		t.Fatalf("failed to unmarshal error response: %v. Body: %s", err, resp.Body.String())
	}

	message, ok := errorResponse["message"].(string)
	if !ok {
		t.Fatalf("expected 'message' field in error response, got none or wrong type. Response: %+v", errorResponse)
	}

	expectedFullErrorMessage := fmt.Sprintf("Failed to process DSL '%s': %s", expectedDSL, expectedErrorMessage)
	if message != expectedFullErrorMessage {
		 t.Errorf("expected error message '%s', got '%s'", expectedFullErrorMessage, message)
	}

	t.Log("SaveMockContent_DSLError completed. Response:", resp.Body.String())
}

// TODO: Implement TestMockContentController_UpdateMockContent_WithDSL_Success
// This test would be structured similarly to SaveMockContent_WithDSL_Success:
// - Setup mocks:
//   - mockProjectSvc.GetProjectBySlugFunc
//   - mockUrlSvc.GetURLByIDFunc (to return an existing URL)
//   - mockFakerSvc.ProcessDSLFunc (to return successful processed data)
//   - mockMcSvc.UpdateMockContentListFunc (to capture updated mock contents)
// - Prepare request payload with dtos.UpdateMockContentUrlDTO which uses []dtos.MockContentUpdateDTO.
// - Make a PATCH request to "/mock/:projectSlug/:urlId".
// - Assertions:
//   - HTTP status OK (or appropriate success status for update).
//   - mockMcSvc.UpdateMockContentListFunc was called with data matching mockFakerSvc output.

// TODO: Implement TestMockContentController_UpdateMockContent_DSLError
// This test would be structured similarly to SaveMockContent_DSLError:
// - Setup mocks:
//   - mockProjectSvc.GetProjectBySlugFunc
//   - mockUrlSvc.GetURLByIDFunc
//   - mockFakerSvc.ProcessDSLFunc (to return an error)
// - Prepare request payload (dtos.UpdateMockContentUrlDTO).
// - Make a PATCH request.
// - Assertions:
//   - HTTP status InternalServerError.
//   - Response body contains the expected error message originating from FakerService.

func TestMockContentController_UpdateMockContent_WithDSL_Success_TODO(t *testing.T) {
    t.Log("TestMockContentController_UpdateMockContent_WithDSL_Success needs to be implemented with proper Gin context and mocking of dependent services, similar to Save tests.")
}

func TestMockContentController_UpdateMockContent_DSLError_TODO(t *testing.T) {
    t.Log("TestMockContentController_UpdateMockContent_DSLError needs to be implemented, similar to Save tests but for PATCH and Update DTOs.")
}
