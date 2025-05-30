package middleware

import (
	"github.com/gin-gonic/gin"
)

// HeaderMiddleware extracts X-Team-Slug and X-Project-Slug from request headers
// and sets them in the Gin context for downstream handlers.
func HeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		teamSlug := c.GetHeader("X-Team-Slug")
		projectSlug := c.GetHeader("X-Project-Slug")

		if teamSlug != "" {
			c.Set("teamSlug", teamSlug)
		}

		if projectSlug != "" {
			c.Set("projectSlug", projectSlug)
		}

		c.Next()
	}
}
