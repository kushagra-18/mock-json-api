package utils

import (
	"strings"
)

// LTrimChar removes all leading instances of a specific character (rune) from a string.
// For example, LTrimChar("/my/path", '/') would return "my/path".
func LTrimChar(s string, charToRemove rune) string {
	// strings.TrimLeft expects a cutset string.
	// So, we convert the rune to a string.
	return strings.TrimLeft(s, string(charToRemove))
}

// Example of another utility if needed:
// func GenerateRandomString(length int) (string, error) {
// 	bytes := make([]byte, length)
// 	if _, err := rand.Read(bytes); err != nil {
// 		return "", err
// 	}
// 	return hex.EncodeToString(bytes), nil
// }
