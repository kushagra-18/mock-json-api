package utils_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"mockapi/utils" // Adjust import path if your module structure is different
)

const (
	testSecretKey     = "testsecret"
	anotherSecretKey  = "anothersecret"
	testUserID        = "user123"
	testTeamID        = "team456"
	shortExpiration   = 1 * time.Millisecond // For testing expiration
	normalExpiration  = 1 * time.Hour
)

func TestGenerateJWTToken(t *testing.T) {
	t.Run("success_case", func(t *testing.T) {
		tokenString, err := utils.GenerateJWTToken(testUserID, testTeamID, testSecretKey, normalExpiration)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)
	})
}

func TestValidateJWTToken(t *testing.T) {
	t.Run("success_case_valid_token", func(t *testing.T) {
		tokenString, err := utils.GenerateJWTToken(testUserID, testTeamID, testSecretKey, normalExpiration)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		token, err := utils.ValidateJWTToken(tokenString, testSecretKey)
		assert.NoError(t, err)
		assert.NotNil(t, token)
		assert.True(t, token.Valid)

		claims, ok := token.Claims.(*utils.AppClaims)
		assert.True(t, ok)
		assert.Equal(t, testUserID, claims.UserID)
		assert.Equal(t, testTeamID, claims.TeamID)
		assert.Equal(t, "mockapi", claims.Issuer) // Default issuer from GenerateJWTToken
	})

	t.Run("error_case_invalid_token_string", func(t *testing.T) {
		_, err := utils.ValidateJWTToken("this.is.not.a.valid.token", testSecretKey)
		assert.Error(t, err)
		// Error message changed after ValidateJWTToken modifications
		assert.Contains(t, err.Error(), "failed to parse token: token contains an invalid number of segments")
	})

	t.Run("error_case_wrong_secret_key", func(t *testing.T) {
		tokenString, err := utils.GenerateJWTToken(testUserID, testTeamID, testSecretKey, normalExpiration)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		_, err = utils.ValidateJWTToken(tokenString, anotherSecretKey)
		assert.Error(t, err)
		// Error from jwt-go library for signature mismatch
		assert.Contains(t, err.Error(), "signature is invalid")
	} )

	t.Run("error_case_expired_token", func(t *testing.T) {
		tokenString, err := utils.GenerateJWTToken(testUserID, testTeamID, testSecretKey, shortExpiration)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		// Wait for token to expire
		time.Sleep(shortExpiration + 5*time.Millisecond) // Sleep a bit longer than expiration

		_, err = utils.ValidateJWTToken(tokenString, testSecretKey)
		assert.Error(t, err)
		// Error from jwt-go for expired token
		var validationError *jwt.ValidationError
		assert.ErrorAs(t, err, &validationError) // Check if it's a jwt.ValidationError
		if validationError != nil {
			assert.True(t, validationError.Errors&jwt.ValidationErrorExpired != 0)
		}
	})

	t.Run("error_case_malformed_token_bad_claims", func(t *testing.T) {
		// Create a token with a different claims type
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(normalExpiration)),
			Issuer:    "anotherissuer",
		})
		malformedTokenString, err := token.SignedString([]byte(testSecretKey))
		assert.NoError(t, err)

		_, err = utils.ValidateJWTToken(malformedTokenString, testSecretKey)
		assert.Error(t, err)
		// This might manifest as a general validation error or a specific claims type assertion error
		// depending on how ParseWithClaims handles it internally when the provided &AppClaims{} doesn't match.
		// The jwt.ParseWithClaims will attempt to unmarshal into AppClaims. If it fails due to type mismatch,
		// it might not error out immediately but result in `token.Valid` being false or claims not being AppClaims type.
		// However, our ValidateJWTToken returns an error if token.Valid is false after parsing.
		// The current error message "failed to parse or validate token" is generic.
		// For a more specific test, one might need to inspect the returned token if err was nil but token.Valid was false.
		// But ValidateJWTToken already wraps this.
		// After changes to ValidateJWTToken, it should now error due to missing UserID.
		assert.Contains(t, err.Error(), "token is missing required custom claim: UserID")
	})
}
