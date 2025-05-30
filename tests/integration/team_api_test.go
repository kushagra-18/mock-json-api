package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"

	"go-gin-gorm-api/internal/handlers"
	"go-gin-gorm-api/internal/models"
	"go-gin-gorm-api/internal/repositories" // For the repository interface
	"go-gin-gorm-api/internal/services"
)

// MockTeamRepository is a local mock implementation of repositories.TeamRepository for this integration test.
type MockTeamRepository struct {
	mock.Mock
}

// Ensure MockTeamRepository implements repositories.TeamRepository
var _ repositories.TeamRepository = &MockTeamRepository{}

func (m *MockTeamRepository) CreateTeam(team *models.Team) error {
	args := m.Called(team)
	return args.Error(0)
}

func (m *MockTeamRepository) GetTeamByID(id uint) (*models.Team, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Team), args.Error(1)
}

func (m *MockTeamRepository) GetTeamBySlug(slug string) (*models.Team, error) {
	args := m.Called(slug)
	// Handle the case where the first return argument (the *models.Team) might be nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Team), args.Error(1)
}

func (m *MockTeamRepository) UpdateTeam(team *models.Team) error {
	args := m.Called(team)
	return args.Error(0)
}

func (m *MockTeamRepository) DeleteTeam(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTeamRepository) GetAllTeams() ([]models.Team, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Team), args.Error(1)
}


// TestTeamAPI_CreateTeam_Success tests the successful creation of a team via the API.
func TestTeamAPI_CreateTeam_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup
	mockRepo := new(MockTeamRepository)
	teamService := services.NewTeamService(mockRepo) // Using the local mock
	teamHandler := handlers.NewTeamHandler(teamService)

	router := gin.New()
	apiV1Group := router.Group("/api/v1")
	// Assuming RegisterTeamRoutes uses the group for /teams, so routes become /api/v1/teams
	handlers.RegisterTeamRoutes(apiV1Group, teamHandler)

	// Payload and Mocking
	teamPayload := `{"name": "Integration Test Team", "slug": "custom-integration-slug"}` // Provided slug for simplicity in this test
	expectedName := "Integration Test Team"
	expectedSlug := "custom-integration-slug"

	// Mock for CreateTeam:
	// Since slug is provided, TeamService will use it directly. No GetTeamBySlug call for collision check.
	mockRepo.On("CreateTeam", mock.MatchedBy(func(team *models.Team) bool {
		return team.Name == expectedName && team.Slug == expectedSlug
	})).Run(func(args mock.Arguments) {
		teamToCreate := args.Get(0).(*models.Team)
		teamToCreate.ID = 1 // Simulate database assigning an ID
		// CreatedAt, UpdatedAt would also be set by DB/GORM
	}).Return(nil).Once()

	// Request
	req, err := http.NewRequest(http.MethodPost, "/api/v1/teams", strings.NewReader(teamPayload))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	// Execution
	router.ServeHTTP(rr, req)

	// Assertions
	assert.Equal(t, http.StatusCreated, rr.Code)

	var responseTeam models.Team
	err = json.Unmarshal(rr.Body.Bytes(), &responseTeam)
	assert.NoError(t, err)

	assert.Equal(t, expectedName, responseTeam.Name)
	assert.Equal(t, expectedSlug, responseTeam.Slug)
	assert.Equal(t, uint(1), responseTeam.ID) // Check if ID set by mock is returned
	// Timestamps would be zero/default if not set in mock and not part of this specific test focus

	mockRepo.AssertExpectations(t)
}

func TestTeamAPI_CreateTeam_SlugGeneration(t *testing.T) {
    gin.SetMode(gin.TestMode)

    mockRepo := new(MockTeamRepository)
    teamService := services.NewTeamService(mockRepo)
    teamHandler := handlers.NewTeamHandler(teamService)

    router := gin.New()
    apiV1Group := router.Group("/api/v1")
    handlers.RegisterTeamRoutes(apiV1Group, teamHandler)

    teamPayload := `{"name": "Integration Test Team Auto Slug"}` // Slug not provided
    expectedName := "Integration Test Team Auto Slug"
    expectedSlugPrefix := "integration-test-team-auto-slug"

    // Mock for GetTeamBySlug (collision check for auto-generated slug)
    // We expect it to be called with the base slug first.
    mockRepo.On("GetTeamBySlug", expectedSlugPrefix).Return(nil, gorm.ErrRecordNotFound).Once()

    // Mock for CreateTeam
    mockRepo.On("CreateTeam", mock.MatchedBy(func(team *models.Team) bool {
        return team.Name == expectedName && team.Slug == expectedSlugPrefix // In this non-collision case
    })).Run(func(args mock.Arguments) {
        teamToCreate := args.Get(0).(*models.Team)
        teamToCreate.ID = 2 // Simulate DB assigning an ID
    }).Return(nil).Once()

    // Request
    req, err := http.NewRequest(http.MethodPost, "/api/v1/teams", strings.NewReader(teamPayload))
    assert.NoError(t, err)
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    // Execution
    router.ServeHTTP(rr, req)

    // Assertions
    assert.Equal(t, http.StatusCreated, rr.Code)

    var responseTeam models.Team
    err = json.Unmarshal(rr.Body.Bytes(), &responseTeam)
    assert.NoError(t, err)

    assert.Equal(t, expectedName, responseTeam.Name)
    assert.Equal(t, expectedSlugPrefix, responseTeam.Slug) // Slug should be the auto-generated one
    assert.Equal(t, uint(2), responseTeam.ID)

    mockRepo.AssertExpectations(t)
}
