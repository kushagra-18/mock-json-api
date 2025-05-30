package dtos

// URLDataDTO is used for transferring URL data, particularly for updates.
// Pointers are used for nullable fields to distinguish between a zero value and a field not being set.
type URLDataDTO struct {
	Description *string `json:"description"`
	Name        *string `json:"name"`
	Requests    *int64  `json:"requests"` // Changed from int to int64 to match model's sql.NullInt64
	Time        *int64  `json:"time"`     // Changed from int to int64 to match model's sql.NullInt64
	Status      *string `json:"status"`   // Added to allow status updates, maps to models.StatusCode
}
