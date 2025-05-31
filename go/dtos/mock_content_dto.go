package dtos

import "mockapi/models" // For models.StatusCode

// MockContentCreateDTO is used for creating a new mock content item.
// Omits ID, CreatedAt, UpdatedAt, DeletedAt, UrlID (set by service).
type MockContentCreateDTO struct {
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	Data        string  `json:"data" binding:"required"` // Assuming data is always required
	Randomness  *int64  `json:"randomness,omitempty"`    // Use omitempty for optional fields with defaults
	Latency     *int64  `json:"latency,omitempty"`
}

// MockContentUpdateDTO is used for updating an existing mock content item.
// ID is essential to identify the item to update. Other fields are optional.
type MockContentUpdateDTO struct {
	ID          *uint   `json:"id"` // Pointer to allow null if it's a new item in an update list, though typically required for updates
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Data        *string `json:"data"`
	Randomness  *int64  `json:"randomness"`
	Latency     *int64  `json:"latency"`
}

// MockContentUrlDTO is used for creating a URL along with its mock contents.
// This corresponds to `com.mock_json.mock_api.dtos.MockContentUrlDto`.
type MockContentUrlDTO struct {
	// URLData contains fields for creating/updating the models.Url itself
	URLData struct {
		Description *string           `json:"description"`
		Name        string            `json:"name" binding:"required"`
		URL         string            `json:"url" binding:"required"` // The path for the URL
		Status      models.StatusCode `json:"status" binding:"required"`
		// Requests and Time are usually not set at creation, but managed by system.
	} `json:"url_data" binding:"required"`

	MockContentList []MockContentCreateDTO `json:"mock_content_list" binding:"required,dive"` // dive validates each element in slice
}

// UpdateMockContentUrlDTO is used for updating a URL's mock contents.
// Typically, this would involve sending the full list of current mock contents for that URL.
// The service layer would then diff or replace them.
type UpdateMockContentUrlDTO struct {
	// List of mock contents. Items with ID are updated, items without ID (or ID=0) might be created.
	// Items missing from the list compared to DB might be deleted by the service.
	MockContentList []MockContentUpdateDTO `json:"mock_content_list" binding:"required,dive"`
}

// GetMockedJSONParamsDTO defines parameters that might be passed for the GetMockedJSON endpoint,
// typically via query params or headers, after base64 decoding.
// This is a conceptual DTO; actual parsing might be more complex.
type GetMockedJSONParamsDTO struct {
	AuthToken     string // From Authorization header
	CustomParam1  string // Example of other decoded params
	Forward       *bool  `json:"forward"`         // From request body after base64 decode
	IsForwardCall *bool  `json:"is_forward_call"` // From request body
}
