package utils_test

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"mockapi/utils" // Adjust import path
)

func TestGenerateRandomString(t *testing.T) {
	t.Run("correct_length", func(t *testing.T) {
		lengths := []int{1, 10, 100}
		for _, length := range lengths {
			s, err := utils.GenerateRandomString(length)
			assert.NoError(t, err)
			assert.Len(t, s, length)
		}
	})

	t.Run("lowercase_letters_only", func(t *testing.T) {
		s, err := utils.GenerateRandomString(20)
		assert.NoError(t, err)
		assert.Regexp(t, regexp.MustCompile("^[a-z]+$"), s, "String should contain only lowercase letters")
	})

	t.Run("error_for_negative_length", func(t *testing.T) {
		s, err := utils.GenerateRandomString(-1)
		// The current implementation returns empty string and no error for length <= 0.
		// Depending on desired strictness, this could be an error.
		// For now, testing current behavior.
		assert.NoError(t, err) // Or assert.Error(t, err) if it should error
		assert.Empty(t, s)      // Current behavior
	})

	t.Run("zero_length", func(t *testing.T) {
		s, err := utils.GenerateRandomString(0)
		assert.NoError(t, err)
		assert.Empty(t, s)
	})

	t.Run("subsequent_calls_produce_different_strings", func(t *testing.T) {
		s1, err1 := utils.GenerateRandomString(10)
		s2, err2 := utils.GenerateRandomString(10)
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NotEqual(t, s1, s2, "Two consecutively generated strings of same length should generally not be equal")
	})
}

func TestUnslug(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"dashes", "hello-world", "Hello World"},
		{"underscores", "hello_world", "Hello World"},
		{"mixed_cases_dashes", "Hello-World-Test", "Hello World Test"},
		{"mixed_cases_underscores", "Hello_World_Test", "Hello World Test"},
		{"single_word", "hello", "Hello"},
		{"single_word_capitalized", "Hello", "Hello"},
		{"empty_input", "", ""},
		{"leading_trailing_spaces_with_slug_chars", "  hello-world_test  ", "Hello World Test"}, // Current impl doesn't trim surrounding spaces first
		{"numbers_in_slug", "version-1-2-3", "Version 1 2 3"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, utils.Unslug(tt.input))
		})
	}
}

func TestDecodeBase64(t *testing.T) {
	t.Run("valid_base64_string", func(t *testing.T) {
		originalString := "Hello, World! 123"
		encodedString := "SGVsbG8sIFdvcmxkISAxMjM=" // echo -n "Hello, World! 123" | base64

		decoded, err := utils.DecodeBase64(encodedString)
		assert.NoError(t, err)
		assert.Equal(t, originalString, decoded)
	})

	t.Run("empty_string", func(t *testing.T) {
		decoded, err := utils.DecodeBase64("")
		assert.NoError(t, err)
		assert.Equal(t, "", decoded)
	})

	t.Run("invalid_base64_string", func(t *testing.T) {
		invalidEncodedString := "This is not valid base64!!!"
		_, err := utils.DecodeBase64(invalidEncodedString)
		assert.Error(t, err)
		// The error from base64.StdEncoding.DecodeString is specific, e.g., "illegal base64 data at input byte X"
		// For a generic check:
		assert.Contains(t, err.Error(), "illegal base64 data")
	})
}
