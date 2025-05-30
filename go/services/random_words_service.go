package services

import (
	"fmt"
	"math/rand"
	"strings"
	// "time" // Only needed for seeding rand if not using Go 1.20+ global seeding
)

// Predefined lists of words for generating random slugs.
// In a real application, these might come from a configuration file or database.
var (
	adjectives = []string{
		"autumn", "hidden", "bitter", "misty", "silent", "empty", "dry", "dark",
		"summer", "icy", "delicate", "quiet", "white", "cool", "spring", "winter",
		"patient", "twilight", "dawn", "crimson", "wispy", "weathered", "blue",
		"billowing", "broken", "cold", "damp", "falling", "frosty", "green",
		"long", "late", "lingering", "bold", "little", "morning", "muddy", "old",
		"red", "rough", "still", "small", "sparkling", "throbbing", "shy",
		"wandering", "withered", "wild", "black", "young", "holy", "solitary",
		"fragrant", "aged", "snowy", "proud", "floral", "restless", "divine",
		"polished", "ancient", "purple", "lively", "nameless", "lucky", "odd",
		"tiny", "soft", "cool", "rapid", "shy", "sweet", "valiant", "warm",
	}
	nouns = []string{
		"waterfall", "river", "breeze", "moon", "rain", "wind", "sea", "morning",
		"snow", "lake", "sunset", "pine", "shadow", "leaf", "dawn", "glitter",
		"forest", "hill", "cloud", "meadow", "sun", "glade", "bird", "brook",
		"butterfly", "bush", "dew", "dust", "field", "fire", "flower", "firefly",
		"feather", "grass", "haze", "mountain", "night", "pond", "darkness",
		"snowflake", "silence", "sound", "sky", "shape", "surf", "thunder",
		"violet", "water", "wildflower", "wave", "water", "resonance", "sun",
		"wood", "dream", "cherry", "tree", "fog", "frost", "voice", "paper",
		"frog", "smoke", "star", "atom", "band", "bar", "base", "block", "boat",
		"term", "credit", "art", "fashion", "truth", "way", "wisdom", "token",
	}
	// List of disallowed slugs (e.g., reserved words, routes)
	// Should be kept lowercase for case-insensitive comparison.
	disallowedProjectSlugs = map[string]bool{
		"new":    true,
		"edit":   true,
		"api":    true,
		"admin":  true,
		"assets": true,
		"static": true,
		"auth":   true,
		"login":  true,
		"logout": true,
		"user":   true,
		"users":  true,
		"team":   true,
		"teams":  true,
		"project":true,
		"projects":true,
		"settings": true,
		"profile": true,
		"dashboard": true,
		"help": true,
		"support": true,
		"status": true,
		"search": true,
		"metrics": true,
		"public": true,
	}
)

// RandomWordsService generates random slugs.
type RandomWordsService struct {
	// No dependencies like DB needed for this version.
	// randSource *rand.Rand // Use if specific seeding or source is needed.
}

// NewRandomWordsService creates a new RandomWordsService.
func NewRandomWordsService() *RandomWordsService {
	// For Go versions before 1.20, seeding was manual:
	// source := rand.NewSource(time.Now().UnixNano())
	// return &RandomWordsService{randSource: rand.New(source)}
	// For Go 1.20+, global rand is auto-seeded and concurrency-safe.
	return &RandomWordsService{}
}

// getRandInt generates a random integer up to max.
// func (s *RandomWordsService) getRandInt(max int) int {
//  if s.randSource != nil {
// 		return s.randSource.Intn(max)
// 	}
// 	return rand.Intn(max)
// }

// GetRandomSlug generates a random, possibly multi-word slug (e.g., "adjective-noun-number")
// and ensures it's not in a disallowed list.
func (s *RandomWordsService) GetRandomSlug() string {
	maxAttempts := 10 // Avoid infinite loops if somehow all generated slugs are disallowed
	for i := 0; i < maxAttempts; i++ {
		adj := adjectives[rand.Intn(len(adjectives))]
		noun := nouns[rand.Intn(len(nouns))]
		num := rand.Intn(9000) + 1000 // Random number between 1000 and 9999

		slug := fmt.Sprintf("%s-%s-%d", strings.ToLower(adj), strings.ToLower(noun), num)

		if !disallowedProjectSlugs[slug] {
			return slug
		}
	}
	// Fallback if too many attempts fail: generate with a timestamp or longer random number
	return fmt.Sprintf("generated-slug-%d", time.Now().UnixNano())
}

// IsSlugDisallowed checks if a given slug is in the disallowed list.
// Primarily for external use if needed, as GetRandomSlug already performs this check.
func (s *RandomWordsService) IsSlugDisallowed(slug string) bool {
	return disallowedProjectSlugs[strings.ToLower(slug)]
}
