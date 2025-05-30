package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"mockapi/models"
	"mockapi/services" // Adjust import path
)

// TestSelectRandomMockContent tests the weighted random selection logic.
func TestSelectRandomMockContent(t *testing.T) {
	// Service does not have DB dependencies for SelectRandomMockContent
	service := services.NewMockContentService(nil) // No DB needed for this method

	t.Run("empty_list", func(t *testing.T) {
		result := service.SelectRandomMockContent([]models.MockContent{})
		assert.Nil(t, result)
	})

	t.Run("list_with_one_item", func(t *testing.T) {
		item := models.MockContent{Name: "OnlyItem", Randomness: 10}
		contents := []models.MockContent{item}
		result := service.SelectRandomMockContent(contents)
		assert.NotNil(t, result)
		assert.Equal(t, "OnlyItem", result.Name)
	})

	t.Run("all_items_zero_randomness", func(t *testing.T) {
		contents := []models.MockContent{
			{Name: "Item1", Randomness: 0},
			{Name: "Item2", Randomness: 0},
			{Name: "Item3", Randomness: 0},
		}
		// Expect one of them to be chosen (uniformly random)
		// To make this test deterministic, we can't easily predict which one.
		// We can check if the returned item is one of the list.
		found := false
		for i := 0; i < 10; i++ { // Run a few times to increase confidence
			result := service.SelectRandomMockContent(contents)
			assert.NotNil(t, result)
			for _, c := range contents {
				if result.Name == c.Name {
					found = true
					break
				}
			}
			assert.True(t, found, "Selected item should be from the original list")
			found = false // Reset for next iteration
		}
	})

	t.Run("one_item_has_all_weight", func(t *testing.T) {
		contents := []models.MockContent{
			{Name: "Item1", Randomness: 0},
			{Name: "Winner", Randomness: 100},
			{Name: "Item3", Randomness: 0},
		}
		// With non-zero total randomness, the one with weight should always be chosen.
		for i := 0; i < 10; i++ { // Run a few times to be sure
			result := service.SelectRandomMockContent(contents)
			assert.NotNil(t, result)
			assert.Equal(t, "Winner", result.Name)
		}
	})

	t.Run("negative_randomness_treated_as_zero", func(t *testing.T) {
		contents := []models.MockContent{
			{Name: "Item1", Randomness: -10},
			{Name: "Winner", Randomness: 1}, // Only one with positive weight
			{Name: "Item3", Randomness: -5},
		}
		result := service.SelectRandomMockContent(contents)
		assert.NotNil(t, result)
		assert.Equal(t, "Winner", result.Name)
	})


	t.Run("probabilistic_distribution_rough_check", func(t *testing.T) {
		// This is a more complex test. For a simple check:
		// If ItemA has weight 90 and ItemB has 10, ItemA should be chosen much more often.
		contents := []models.MockContent{
			{Name: "RareItem", Randomness: 10},   // 10% chance (approx)
			{Name: "CommonItem", Randomness: 90}, // 90% chance (approx)
		}
		counts := map[string]int{"RareItem": 0, "CommonItem": 0}
		numRuns := 1000
		for i := 0; i < numRuns; i++ {
			result := service.SelectRandomMockContent(contents)
			if result != nil {
				counts[result.Name]++
			}
		}
		// Expected: CommonItem count should be significantly higher than RareItem.
		// Allow some variance, e.g., CommonItem > 700, RareItem < 300 for numRuns=1000
		// These are not strict statistical bounds but a sanity check.
		// t.Logf("Counts: RareItem=%d, CommonItem=%d", counts["RareItem"], counts["CommonItem"])
		assert.True(t, counts["CommonItem"] > counts["RareItem"], "CommonItem should be selected more often")
		assert.True(t, counts["CommonItem"] > (numRuns/2), "CommonItem should be selected more than half the time")
		assert.True(t, counts["RareItem"] < (numRuns/2), "RareItem should be selected less than half the time")
		assert.Equal(t, numRuns, counts["RareItem"]+counts["CommonItem"], "Total selections should match number of runs")
	})
}

// TestSimulateLatency tests the time.Sleep functionality.
func TestSimulateLatency(t *testing.T) {
	service := services.NewMockContentService(nil)

	t.Run("positive_latency", func(t *testing.T) {
		latencyMillis := int64(20) // Small latency
		startTime := time.Now()
		service.SimulateLatency(latencyMillis)
		duration := time.Since(startTime)

		// Check if duration is approximately latencyMillis.
		// Allow for some overhead/inaccuracy in sleep and measurement.
		// Lower bound: at least the latency. Upper bound: latency + some buffer (e.g., 15ms).
		minExpectedDuration := time.Duration(latencyMillis) * time.Millisecond
		maxExpectedDuration := time.Duration(latencyMillis+15) * time.Millisecond

		assert.True(t, duration >= minExpectedDuration, "Duration (%v) should be at least the specified latency (%v)", duration, minExpectedDuration)
		assert.True(t, duration <= maxExpectedDuration, "Duration (%v) should not significantly exceed latency (%v)", duration, maxExpectedDuration)
	})

	t.Run("zero_latency", func(t *testing.T) {
		latencyMillis := int64(0)
		startTime := time.Now()
		service.SimulateLatency(latencyMillis)
		duration := time.Since(startTime)
		// Should be very quick, less than a few milliseconds
		assert.True(t, duration < (5*time.Millisecond), "Duration for zero latency should be minimal, got %v", duration)
	})

	t.Run("negative_latency", func(t *testing.T) {
		latencyMillis := int64(-10)
		startTime := time.Now()
		service.SimulateLatency(latencyMillis) // Should behave like zero latency
		duration := time.Since(startTime)
		assert.True(t, duration < (5*time.Millisecond), "Duration for negative latency should be minimal, got %v", duration)
	})
}
