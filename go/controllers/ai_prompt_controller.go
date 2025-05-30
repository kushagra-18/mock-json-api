package controllers

import (
	"net/http"
	// "log" // For potential detailed error logging

	"github.com/gin-gonic/gin"
	// "google.golang.org/genai" // To check for genai.APIError, if needed

	"mockapi/dtos"
	"mockapi/services"
	"mockapi/utils" // For response helpers
)

// AIPromptController handles API endpoints related to AI prompting.
type AIPromptController struct {
	aiPromptService services.AIPromptServiceInterface // Use the interface type
}

// NewAIPromptController creates a new AIPromptController.
func NewAIPromptController(service services.AIPromptServiceInterface) *AIPromptController { // Accept interface type
	if service == nil {
		// Or panic, depending on desired strictness for dependency injection
		panic("AIPromptServiceInterface cannot be nil in NewAIPromptController")
	}
	return &AIPromptController{aiPromptService: service}
}

// HandleAIPrompt is the Gin handler for AI prompt requests.
func (controller *AIPromptController) HandleAIPrompt(c *gin.Context) {
	var reqDTO dtos.AIPromptRequestDTO

	// Bind and validate the request JSON to the DTO.
	if err := c.ShouldBindJSON(&reqDTO); err != nil {
		// Gin's binding:"required" should catch empty prompts if they are truly empty.
		// If prompt is "   ", it might pass "required" but fail "min=1" after trim.
		// Or, if only "required", "   " might pass. The DTO has `binding:"required,min=1"`.
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	// The DTO binding `min=1` should make this redundant if Gin trims whitespace by default for validation.
	// However, an explicit check after binding can be a safeguard if behavior is uncertain or changes.
	// if strings.TrimSpace(reqDTO.Prompt) == "" {
	// 	 utils.ErrorResponse(c, http.StatusBadRequest, "Prompt cannot be empty.")
	// 	 return
	// }


	// Call the service method, passing c.Request.Context() for context propagation.
	aiResponse, err := controller.aiPromptService.GetAIResponse(c.Request.Context(), reqDTO.Prompt)
	if err != nil {
		// Log the detailed error for server-side observability
		// log.Printf("Error from AIPromptService.GetAIResponse: %v", err)

		// Consider checking the error type for more specific client responses.
		// Example:
		// var apiError *genai.APIError // Assuming genai.APIError is a concrete type or interface
		// if errors.As(err, &apiError) {
		//     // Handle specific API errors, e.g., rate limits, auth issues from Gemini
		//     // This depends on the structure of errors returned by the genai package
		//     utils.ErrorResponse(c, http.StatusServiceUnavailable, "AI service error: "+apiError.Message)
		// } else {
		//     utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get AI response.")
		// }
		// For now, a generic error message is returned.
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get AI response: "+err.Error())
		return
	}

	// Return the successful response.
	utils.SuccessResponse(c, http.StatusOK, aiResponse)
}
