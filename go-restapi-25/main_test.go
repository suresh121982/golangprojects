// main_test.go
package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Helper function to generate a test JWT token
func generateTestToken(username string) string {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(jwtKey)
	return tokenString
}

func TestLoginHandler(t *testing.T) {
	// Define test cases
	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{
			name:       "Valid credentials",
			body:       `{"username":"admin","password":"admin"}`,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid credentials",
			body:       `{"username":"user","password":"pass"}`,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Missing credentials",
			body:       `{}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid JSON",
			body:       `{"username":"admin","password"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Empty body",
			body:       ``,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/login", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json") // Set header to JSON
			w := httptest.NewRecorder()
			loginHandler(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %v, want %v", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}

func TestAuthMiddleware(t *testing.T) {
	// Set up a test server with the authMiddleware
	handler := authMiddleware(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Success"))
	})

	// Valid token generated with 'admin' user
	validToken := generateTestToken("admin")

	// Expired token test setup
	expiredToken := generateTestToken("admin")
	expiredToken = strings.Replace(expiredToken, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHBpcmF0aW9uVGltZSI6IjE2MjM5NTE5MjQiLCJzdWIiOiJhZG1pbiJ9.fW9bRV4K5U9xZ5CmlB_wDsC1A7o0TY4NCjJbUPKgPkk", "invalid-part", -1) // Simulate expiration

	// Test cases for authMiddleware
	tests := []struct {
		name       string
		token      string
		wantStatus int
	}{
		{
			name:       "Valid token",
			token:      validToken,
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid token",
			token:      "invalid-token",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Expired token",
			token:      expiredToken,
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "No token",
			token:      "",
			wantStatus: http.StatusUnauthorized,
		},
		{
			name:       "Valid token but wrong claim",
			token:      generateTestToken("wrong-user"),
			wantStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/restricted", nil)
			if tt.token != "" {
				req.AddCookie(&http.Cookie{Name: "jwt", Value: tt.token})
			}
			w := httptest.NewRecorder()
			handler.ServeHTTP(w, req)

			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %v, want %v", resp.StatusCode, tt.wantStatus)
			}
		})
	}
}
