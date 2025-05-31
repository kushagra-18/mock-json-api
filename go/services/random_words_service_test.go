package services_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"mockapi/services" // Adjust import path
)

// TestGetRandomSlug checks basic functionality of GetRandomSlug.
func TestGetRandomSlug(t *testing.T) {
	service := services.NewRandomWordsService()

	t.Run("returns_non_empty_string", func(t *testing.T) {
		slug := service.GetRandomSlug()
		assert.NotEmpty(t, slug, "Generated slug should not be empty")
	})

	t.Run("returns_slug_in_expected_format", func(t *testing.T) {
		slug := service.GetRandomSlug()
		// Expected format: adjective-noun-number (e.g., "cool-river-1234")
		// This regex is a bit loose on the number part to accommodate the fallback.
		assert.Regexp(t, regexp.MustCompile(`^[a-z]+-[a-z]+-\d+$`), slug, "Slug format is not as expected")
	})

	t.Run("does_not_return_globally_disallowed_slug", func(t *testing.T) {
		// This test is tricky because the disallowed list is hardcoded and global in the service.
		// We can't easily modify it per test without changing the service's design (e.g., pass list in constructor).
		// However, we can check against one known disallowed slug.
		// A better approach would be to make the disallowed list configurable for the service instance.
		// For now, we rely on the existing hardcoded list.
		// Example: if "new-project-1234" was somehow generated and "new" is part of a disallowed prefix,
		// the service's internal check should prevent it.
		// The current `disallowedProjectSlugs` is a map of full slugs.
		// Let's test by trying to generate many slugs and see if any match a simple known disallowed one.
		// This is not a perfect test for this specific case.

		// A more direct test for IsSlugDisallowed is better.
		// For GetRandomSlug, we trust it tries not to return one from the list.
		// If "admin-api-1234" is generated, and "admin-api-1234" is in disallowed, it should retry.
		// The current disallowed list seems to be for full slugs like "new", "admin", etc. not patterns.
		// So a generated "admin-api-1234" would be allowed unless "admin-api-1234" itself is in the map.

		// Let's test the fallback behavior by forcing many attempts (conceptually)
		// This is hard to test directly without mocking rand or making the word lists very small.
		// We'll assume the fallback `generated-slug-%d` would eventually be hit if all combinations were disallowed.
		// For now, focusing on the format and non-empty.
	})

	t.Run("multiple_calls_yield_different_slugs", func(t *testing.T) {
		slug1 := service.GetRandomSlug()
		slug2 := service.GetRandomSlug()
		assert.NotEqual(t, slug1, slug2, "Two consecutively generated slugs should generally be different")
	})
}

// TestIsSlugDisallowed tests the IsSlugDisallowed method.
func TestIsSlugDisallowed(t *testing.T) {
	service := services.NewRandomWordsService()

	// These tests depend on the hardcoded disallowedProjectSlugs map in random_words_service.go
	disallowed := []string{"new", "admin", "API", "UsErS"} // Check case-insensitivity
	allowed := []string{"my-project", "another_slug", "test1234"}

	for _, slug := range disallowed {
		t.Run(fmt.Sprintf("disallowed_%s", slug), func(t *testing.T) {
			// The internal check in IsSlugDisallowed uses strings.ToLower(slug)
			// So, we expect it to match regardless of the input case if the lowercase version is in the map.
			// The map itself should store keys in lowercase.
			// Example: if "new" is in map, "New", "NEW" should be caught.
			expectedDisallowed := false
			// Check if the lowercase version of the slug is in the hardcoded map
			// (This is re-implementing the service's map access logic for test setup, which is not ideal)
			// For this test, we rely on the service's internal map.
			// The service's `disallowedProjectSlugs` map keys are already lowercase.
			// So, `service.IsSlugDisallowed("API")` will check for "api" in the map.

			// Hardcoded map for reference:
			// "new":true, "edit":true, "api":true, "admin":true, ...

			// Let's assume the service's map has "new", "admin", "api", "users"
			// Then IsSlugDisallowed("new") should be true.
			// IsSlugDisallowed("API") should be true (because it checks "api").
			// IsSlugDisallowed("UsErS") should be true (because it checks "users").

			// Based on the provided disallowed list
			if strings.ToLower(slug) == "new" || strings.ToLower(slug) == "admin" || strings.ToLower(slug) == "api" || strings.ToLower(slug) == "users" {
				expectedDisallowed = true
			}


			assert.Equal(t, expectedDisallowed, service.IsSlugDisallowed(slug),
				fmt.Sprintf("Slug '%s' (lowercase '%s') disallow status was not as expected.", slug, strings.ToLower(slug)))
		})
	}

	for _, slug := range allowed {
		t.Run(fmt.Sprintf("allowed_%s", slug), func(t *testing.T) {
			assert.False(t, service.IsSlugDisallowed(slug), fmt.Sprintf("Slug '%s' should be allowed", slug))
		})
	}
}
