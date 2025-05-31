package utils

import "github.com/gin-gonic/gin"

// SuccessResponse sends a standardized success JSON response.
func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, gin.H{
		"status": "success",
		"data":   data,
	})
}

// ErrorResponse sends a standardized error JSON response.
func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{
		"status":  "error",
		"message": message,
	})
}

// CustomResponse allows for more flexible JSON responses.
func CustomResponse(c *gin.Context, statusCode int, payload gin.H) {
	c.JSON(statusCode, payload)
}
