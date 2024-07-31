// utils.go

package utils

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var (
	// SecretKey used to sign JWT tokens (replace with your secret key)
	SecretKey = []byte("your-secret-key")
)

// Claims defines the custom claims for JWT token
type Claims struct {
	UserID int `json:"userID"`
	jwt.StandardClaims
}

// GenerateToken generates a new JWT token for a given user ID
func GenerateToken(userID int) (string, error) {
	// Define token claims
	claims := Claims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
			IssuedAt:  time.Now().Unix(),
			Subject:   fmt.Sprintf("%d", userID),
		},
	}

	// Create a new token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate the token string
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the user ID if valid
func ValidateToken(tokenString string) (int, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		return 0, err
	}

	// Check if the token is valid
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	// Extract and return the user ID from claims
	return claims.UserID, nil
}
