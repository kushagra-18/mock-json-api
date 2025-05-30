package services

import (
	"errors"
	"go-gin-gorm-api/internal/models"
	"go-gin-gorm-api/internal/repositories" // To define the interface for the mock
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockTeamRepository is a mock implementation of TeamRepository for testing.
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

// --- Test Functions ---

// TestTeamService_CreateTeam_SlugGeneration tests the slug generation logic in CreateTeam.
func TestTeamService_CreateTeam_SlugGeneration(t *testing.T) {
	mockRepo := new(MockTeamRepository)
	// teamService := NewTeamService(mockRepo)
	// For testing private/unexported generateSlug, we need to use the actual service
	// or make generateSlug part of an interface/public if complex enough.
	// The current teamService.generateSlug is unexported.
	// To test the full CreateTeam behavior including its internal slug generation,
	// we have to modify the service to make slug generation testable,
	// or accept that we are testing the public CreateTeam method which uses it.

	// The current team_service.go has generateSlug as a package-level unexported function.
	// We can test the behavior of CreateTeam which calls it.
	// The service also has its own generateSlugRandomized for collisions.
	// This makes testing the exact slug tricky without also mocking GetTeamBySlug for collision checks.

	teamServiceInstance := NewTeamService(mockRepo).(*teamService) // Cast to access generateSlugRandomized if needed, or just use public methods

	t.Run("Slug provided by user", func(t *testing.T) {
		inputTeam := &models.Team{Name: "Test Team With Provided Slug", Slug: "custom-slug"}
		// Expect CreateTeam to be called with the exact team object (or one with the same slug)
		mockRepo.On("CreateTeam", mock.MatchedBy(func(team *models.Team) bool {
			return team.Slug == "custom-slug" && team.Name == "Test Team With Provided Slug"
		})).Return(nil).Once()

		err := teamServiceInstance.CreateTeam(inputTeam)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Slug not provided, generated from name", func(t *testing.T) {
		inputTeam := &models.Team{Name: "Another Team Name"}
		expectedSlug := "another-team-name" // Based on current generateSlug logic

		// Mock GetTeamBySlug to return "not found" for the generated slug, so no collision.
		mockRepo.On("GetTeamBySlug", expectedSlug).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("CreateTeam", mock.MatchedBy(func(team *models.Team) bool {
			return team.Slug == expectedSlug && team.Name == "Another Team Name"
		})).Return(nil).Once()

		err := teamServiceInstance.CreateTeam(inputTeam)
		assert.NoError(t, err)
		assert.Equal(t, expectedSlug, inputTeam.Slug)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Slug generation with collision, then success", func(t *testing.T) {
		inputTeam := &models.Team{Name: "Existing Name For Slug"}
		baseGeneratedSlug := "existing-name-for-slug"

		// First call to GetTeamBySlug (for baseGeneratedSlug) returns an existing team (collision)
		mockRepo.On("GetTeamBySlug", baseGeneratedSlug).Return(&models.Team{ID: 1, Slug: baseGeneratedSlug}, nil).Once()

		// Second call to GetTeamBySlug (for baseGeneratedSlug + random suffix) returns "not found"
		mockRepo.On("GetTeamBySlug", mock.MatchedBy(func(slug string) bool {
			return strings.HasPrefix(slug, baseGeneratedSlug+"-") && len(slug) > len(baseGeneratedSlug+"-")
		})).Return(nil, gorm.ErrRecordNotFound).Once()

		// CreateTeam should be called with the randomized slug
		mockRepo.On("CreateTeam", mock.MatchedBy(func(team *models.Team) bool {
			return team.Name == "Existing Name For Slug" &&
				strings.HasPrefix(team.Slug, baseGeneratedSlug+"-") &&
				len(team.Slug) > len(baseGeneratedSlug+"-")
		})).Return(nil).Once()

		err := teamServiceInstance.CreateTeam(inputTeam)
		assert.NoError(t, err)
		assert.True(t, strings.HasPrefix(inputTeam.Slug, baseGeneratedSlug+"-"), "Slug should have a random suffix")
		assert.True(t, len(inputTeam.Slug) > len(baseGeneratedSlug+"-"), "Slug should be longer than base + hyphen")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Slug generation with multiple collisions, then success", func(t *testing.T) {
		inputTeam := &models.Team{Name: "Highly Colliding Name"}
		baseSlug := "highly-colliding-name"

		// Simulate 2 collisions then success
		mockRepo.On("GetTeamBySlug", baseSlug).Return(&models.Team{ID: 1, Slug: baseSlug}, nil).Once() // 1st collision
		mockRepo.On("GetTeamBySlug", mock.MatchedBy(func(s string) bool { return strings.HasPrefix(s, baseSlug+"-") })).Return(&models.Team{ID: 2}, nil).Once() // 2nd collision (on first random attempt)
		mockRepo.On("GetTeamBySlug", mock.MatchedBy(func(s string) bool { return strings.HasPrefix(s, baseSlug+"-") })).Return(nil, gorm.ErrRecordNotFound).Once() // 3rd attempt (2nd random) is unique

		mockRepo.On("CreateTeam", mock.MatchedBy(func(team *models.Team) bool {
			return team.Name == "Highly Colliding Name" && strings.HasPrefix(team.Slug, baseSlug+"-")
		})).Return(nil).Once()

		err := teamServiceInstance.CreateTeam(inputTeam)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})


	t.Run("CreateTeam with empty name (slug should not be generated from empty name)", func(t *testing.T) {
		inputTeam := &models.Team{Name: ""} // Slug is also empty
		// If name is empty, slug remains empty. Repo should handle this (e.g. DB constraint) or service should validate.
		// Current service logic: if team.Slug == "" && team.Name != "", so if name is empty, slug isn't auto-generated.
		// The repository's CreateTeam would be called with an empty slug.
		// This might be an error case depending on DB constraints (e.g. NOT NULL on slug).
		// For this test, we verify that CreateTeam is called with the slug as is (empty).
		mockRepo.On("CreateTeam", mock.MatchedBy(func(team *models.Team) bool {
			return team.Name == "" && team.Slug == ""
		})).Return(errors.New("database constraint: slug cannot be empty")).Once() // Simulate DB error

		err := teamServiceInstance.CreateTeam(inputTeam)
		assert.Error(t, err) // Expecting an error because the mock repo call will return one
		assert.Equal(t, "", inputTeam.Slug) // Slug should remain empty
		mockRepo.AssertExpectations(t)
	})
}


func TestTeamService_UpdateTeam_SlugGeneration(t *testing.T) {
	mockRepo := new(MockTeamRepository)
	teamServiceInstance := NewTeamService(mockRepo).(*teamService)

	t.Run("Update with slug provided, name change does not affect slug", func(t *testing.T) {
		originalTeam := &models.Team{ID: 1, Name: "Original Name", Slug: "original-slug", UpdatedAt: time.Now().Add(-time.Hour)}
		updatePayload := &models.Team{ID: 1, Name: "New Name", Slug: "provided-slug-on-update"} // User provides a new slug

		// Mock GetTeamBySlug for the new provided slug to ensure it's unique (if service checks this, current one doesn't explicitly for update)
		// The service's UpdateTeam currently does: if team.Slug == "" && team.Name != "" { team.Slug = generateSlug(team.Name) }
		// So if a slug is provided in teamUpdates, it's used directly.
		mockRepo.On("UpdateTeam", mock.MatchedBy(func(team *models.Team) bool {
			return team.ID == 1 && team.Name == "New Name" && team.Slug == "provided-slug-on-update"
		})).Return(nil).Once()

		err := teamServiceInstance.UpdateTeam(updatePayload)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update with name change, slug is empty, should generate new slug", func(t *testing.T) {
		// Assume original slug was 'old-name-slug'
		updatePayload := &models.Team{ID: 1, Name: "Updated Name Only", Slug: ""} // Slug is empty in payload
		expectedNewSlug := "updated-name-only"

		// Mock GetTeamBySlug for collision check on the new expected slug
		mockRepo.On("GetTeamBySlug", expectedNewSlug).Return(nil, gorm.ErrRecordNotFound).Once()
		mockRepo.On("UpdateTeam", mock.MatchedBy(func(team *models.Team) bool {
			return team.ID == 1 && team.Name == "Updated Name Only" && team.Slug == expectedNewSlug
		})).Return(nil).Once()

		err := teamServiceInstance.UpdateTeam(updatePayload)
		assert.NoError(t, err)
		assert.Equal(t, expectedNewSlug, updatePayload.Slug)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update with name change, slug is NOT empty, should NOT regenerate slug", func(t *testing.T) {
		updatePayload := &models.Team{ID: 1, Name: "Name Changed Again", Slug: "keep-this-slug"}

		// The service's UpdateTeam: if team.Slug == "" && team.Name != "" ...
		// Since slug is provided, it will not be regenerated.
		mockRepo.On("UpdateTeam", mock.MatchedBy(func(team *models.Team) bool {
			return team.ID == 1 && team.Name == "Name Changed Again" && team.Slug == "keep-this-slug"
		})).Return(nil).Once()

		err := teamServiceInstance.UpdateTeam(updatePayload)
		assert.NoError(t, err)
		assert.Equal(t, "keep-this-slug", updatePayload.Slug) // Slug should be unchanged
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update, slug becomes empty due to name change, collision on new slug", func(t *testing.T) {
		updatePayload := &models.Team{ID: 1, Name: "Colliding Update Name", Slug: ""}
		baseGeneratedSlug := "colliding-update-name"

		mockRepo.On("GetTeamBySlug", baseGeneratedSlug).Return(&models.Team{ID: 2, Slug: baseGeneratedSlug}, nil).Once() // Collision
		mockRepo.On("GetTeamBySlug", mock.MatchedBy(func(s string) bool { return strings.HasPrefix(s, baseGeneratedSlug+"-")})).Return(nil, gorm.ErrRecordNotFound).Once() // Unique random

		mockRepo.On("UpdateTeam", mock.MatchedBy(func(team *models.Team) bool {
			return team.ID == 1 && team.Name == "Colliding Update Name" && strings.HasPrefix(team.Slug, baseGeneratedSlug+"-")
		})).Return(nil).Once()

		err := teamServiceInstance.UpdateTeam(updatePayload)
		assert.NoError(t, err)
		assert.True(t, strings.HasPrefix(updatePayload.Slug, baseGeneratedSlug+"-"))
		mockRepo.AssertExpectations(t)
	})
}

// Additional tests for GetTeamByID, GetTeamBySlug, GetAllTeams, DeleteTeam can be added here.
// They are simpler as they mostly pass through to the repository.

func TestTeamService_GetTeamByID(t *testing.T) {
    mockRepo := new(MockTeamRepository)
    service := NewTeamService(mockRepo)
    expectedTeam := &models.Team{ID: 1, Name: "Test Team", Slug: "test-team"}

    mockRepo.On("GetTeamByID", uint(1)).Return(expectedTeam, nil).Once()

    team, err := service.GetTeamByID(1)
    assert.NoError(t, err)
    assert.NotNil(t, team)
    assert.Equal(t, expectedTeam.Name, team.Name)
    mockRepo.AssertExpectations(t)
}

func TestTeamService_GetTeamByID_NotFound(t *testing.T) {
    mockRepo := new(MockTeamRepository)
    service := NewTeamService(mockRepo)

    mockRepo.On("GetTeamByID", uint(1)).Return(nil, gorm.ErrRecordNotFound).Once()

    team, err := service.GetTeamByID(1)
    assert.Error(t, err)
    assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
    assert.Nil(t, team)
    mockRepo.AssertExpectations(t)
}
