package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"mockapi/config" // Assuming config contains the Node.js service URL
)

// FakerService handles communication with the Node.js Faker DSL processing service.
type FakerService struct {
	NodeJSBaseURL string       // e.g., "http://localhost:3001"
	HTTPClient    *http.Client // For making HTTP requests
}

// NewFakerService creates a new FakerService.
// The nodeJSServiceURL should be the base URL of the Node.js service (e.g., "http://localhost:3001").
func NewFakerService(cfg config.Config) *FakerService {
	return &FakerService{
		NodeJSBaseURL: cfg.NodeJSFakerServiceURL, // Expect this to be in the config
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second, // Sensible default timeout
		},
	}
}

// ProcessDSLRequest is the payload sent to the Node.js service.
type ProcessDSLRequest struct {
	DSL string `json:"dsl"`
}

// ProcessDSL sends the DSL string to the Node.js service and returns the processed data string.
// The returned value from Node.js can be of any type (string, number, array, object).
// We will receive it as json.RawMessage to handle this flexibility and then marshal it back to a string
// to be stored in the `Data` field of `MockContent`.
func (fs *FakerService) ProcessDSL(dslString string) (string, error) {
	requestPayload := ProcessDSLRequest{DSL: dslString}
	payloadBytes, err := json.Marshal(requestPayload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal DSL request payload: %w", err)
	}

	req, err := http.NewRequest("POST", fs.NodeJSBaseURL+"/process-dsl", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", fmt.Errorf("failed to create request to Node.js service: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := fs.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Node.js service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Node.js service returned non-OK status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// The response from Node.js /process-dsl is the raw processed data (string, number, array, object)
	// We need to store this as a single string in our Go model's Data field.
	// So, we read the raw JSON response.
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body from Node.js service: %w", err)
	}

    // The responseBody is the JSON representation of the data.
    // If the DSL was, for example, `{{name.firstName}}`, responseBody would be `"John"`.
    // If it was `{{number.int}}*2`, responseBody would be `[12,34]`.
    // This is already in the correct string format to be stored in the Data field.
	return string(responseBody), nil
}
