package handlers

import "go-gin-gorm-api/internal/models"

// MockContentUrlDto is used for creating a URL and its associated mock content.
type MockContentUrlDto struct {
	UrlData         models.Url           `json:"url_data"`
	MockContentList []models.MockContent `json:"mock_content_list"`
}

// Other DTOs can be added here as needed.
