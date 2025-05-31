package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"mockapi/utils" // Assuming module name is mockapi
)

// JWTAuthMiddleware creates a gin.HandlerFunc for JWT authentication.
// It prioritizes token extraction from the "token" query parameter,
// then checks the "Authorization" header.
func JWTAuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tokenString string

		// 1. Try to get token from "token" query parameter
		tokenString = c.Query("token")

		// 2. If not in query, try "Authorization" header
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
					tokenString = parts[1]
				} else {
					utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid Authorization header format. Expected 'Bearer <token>'.")
					c.Abort()
					return
				}
			}
		}

		// 3. Check if token is missing
		if tokenString == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authentication token is required.")
			c.Abort()
			return
		}

		// 4. Validate the token
		token, err := utils.ValidateJWTToken(tokenString, secretKey)
		if err != nil {
			// Check for specific JWT errors to provide more context if needed
			if e, ok := err.(*jwt.ValidationError); ok {
				if e.Errors&jwt.ValidationErrorMalformed != 0 {
					utils.ErrorResponse(c, http.StatusUnauthorized, "Malformed token.")
				} else if e.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
					utils.ErrorResponse(c, http.StatusUnauthorized, "Token is expired or not yet valid.")
				} else {
					utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token: "+err.Error())
				}
			} else {
				utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token: "+err.Error())
			}
			c.Abort()
			return
		}

		// 5. Check if token is valid and extract claims
		if claims, ok := token.Claims.(*utils.AppClaims); ok && token.Valid {
			// Set claims into context for downstream handlers
			c.Set("userID", claims.UserID)
			c.Set("teamID", claims.TeamID)
			// Can also set the full claims struct if needed
			// c.Set("claims", claims)
		} else {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid token claims.")
			c.Abort()
			return
		}

		// 6. Call c.Next()
		c.Next()
	}
}
