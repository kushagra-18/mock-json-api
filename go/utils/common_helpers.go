package utils

// StringPointerToString safely dereferences a string pointer or returns an empty string if nil.
func StringPointerToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Int64PointerToInt64 safely dereferences an int64 pointer or returns 0 if nil.
// Useful for optional fields with defaults.
func Int64PointerToInt64(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

// BoolPointer returns a pointer to a bool value.
// Useful for optional boolean fields in DTOs or models where you need to distinguish between false and not set.
func BoolPointer(b bool) *bool {
	return &b
}

// UintPointer returns a pointer to a uint value.
func UintPointer(u uint) *uint {
	return &u
}
