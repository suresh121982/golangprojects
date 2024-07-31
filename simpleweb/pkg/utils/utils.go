package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
)

// GenerateRandomString generates a random string of specified length
func GenerateRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// HashPassword hashes a password using SHA-256
func HashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return base64.URLEncoding.EncodeToString(hash[:])
}

// JSONResponse formats data as JSON and writes to ResponseWriter
func JSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Println("Error encoding JSON response:", err)
	}
}
