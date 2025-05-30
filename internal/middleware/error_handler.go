package middleware

import (
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ErrorHandlerMiddleware provides centralized error handling and panic recovery.
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Panic Recovery
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Panic recovered: %v", r)
				// Check if the response has already been written
				if !c.Writer.Written() {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error - panic"})
				}
				c.Abort() // Ensure no further handlers are called
			}
		}()

		c.Next() // Process request by next handlers

		// After c.Next(), check for errors set by handlers
		if len(c.Errors) > 0 {
			// We only handle the first error in the list for simplicity,
			// as Gin typically aborts on the first c.Error() or c.AbortWithError().
			ginErr := c.Errors[0]
			log.Printf("Error caught by middleware: %v, Type: %T, Meta: %v", ginErr.Err, ginErr.Err, ginErr.Meta)


			// Check if the response has already been written by a handler
			if c.Writer.Written() {
				// If response is already written, we can't send a new JSON error response.
				// We just log it. This might happen if a handler writes a success response
				// then encounters an error and calls c.Error().
				log.Printf("Error occurred after response was written: %v", ginErr.Err)
				return
			}

			// Determine the type of error and respond accordingly
			if errors.Is(ginErr.Err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Resource not found"})
			} else {
				// Check for custom validation errors if you have them.
				// Example:
				// var valError *mycustom.ValidationError
				// if errors.As(ginErr.Err, &valError) {
				// 	c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": valError.Error()})
				// } else
				// For now, general internal server error for non-specific errors
				// If a handler explicitly set a status code before erroring (e.g. c.Status(400); c.Error(err)),
				// Gin might use that. But here we are overriding based on error type.
				// If the error message is safe to expose, use it, otherwise generic.
				// For now, using a generic message for 500.
				// If ginErr.Meta has specific status code, could use that too.

				// A simple approach: if the error message is from a known "bad request" scenario, show it.
				// This is a simplification; a more robust system would use custom error types.
				if ginErr.Type == gin.ErrorTypeBind { // Gin's binding errors
					c.JSON(http.StatusBadRequest, gin.H{"error": "Request binding error: " + ginErr.Err.Error()})
				} else {
					// Fallback for other errors
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				}
			}
			// c.Abort() // Not strictly necessary here if c.JSON is the last call and we return
			return // Stop further processing for this request after handling the error
		}

		// Optional: Fallback for unhandled cases if c.Errors is empty but status suggests an issue.
		// This is less common if handlers correctly use c.Error() or c.AbortWithError().
		// status := c.Writer.Status()
		// if status >= 400 && !c.Writer.Written() {
		// 	c.JSON(status, gin.H{"error": "An unspecified error occurred"})
		// }
	}
}
