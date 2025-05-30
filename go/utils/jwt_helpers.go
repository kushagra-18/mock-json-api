package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AppClaims defines the custom claims for JWT.
type AppClaims struct {
	UserID string `json:"userID"`
	TeamID string `json:"teamID"`
	jwt.RegisteredClaims
}

// GenerateJWTToken creates a new JWT token with custom claims.
func GenerateJWTToken(userID string, teamID string, secretKey string, expirationTime time.Duration) (string, error) {
	// Create the claims
	claims := AppClaims{
		UserID: userID,
		TeamID: teamID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "mockapi", // Example issuer
			Subject:   userID,    // Example subject
		},
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateJWTToken parses and validates a JWT token string.
// It returns the parsed token if valid, or an error otherwise.
func ValidateJWTToken(tokenString string, secretKey string) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AppClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse or validate token: %w", err)
	}

	return token, nil
}
