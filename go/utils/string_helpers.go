package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"unicode"
)

// GenerateRandomString generates a cryptographically secure random string of lowercase letters.
func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", nil // Or an error, depending on desired behavior for non-positive length
	}
	const letters = "abcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

// Unslug converts a slug (e.g., "hello-world" or "hello_world")
// to a human-readable string (e.g., "Hello World").
func Unslug(slug string) string {
	if slug == "" {
		return ""
	}
	// Replace hyphens and underscores with spaces
	s := strings.ReplaceAll(slug, "-", " ")
	s = strings.ReplaceAll(s, "_", " ")

	words := strings.Fields(s)
	for i, word := range words {
		if len(word) > 0 {
			runes := []rune(word)
			runes[0] = unicode.ToUpper(runes[0])
			words[i] = string(runes)
		}
	}
	return strings.Join(words, " ")
}

// DecodeBase64 decodes a Base64 encoded string.
func DecodeBase64(encodedString string) (string, error) {
	decodedBytes, err := base64.StdEncoding.DecodeString(encodedString)
	if err != nil {
		return "", err
	}
	return string(decodedBytes), nil
}
