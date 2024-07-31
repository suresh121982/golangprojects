// handlers.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// LoginHandler handles user login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Retrieve username and password from form data
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	// Find user by username
	user := FindUserByUsername(username)
	if user == nil || user.Password != password {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Simulate session management (usually use cookies or tokens)
	sessionID := fmt.Sprintf("session-%d", user.ID)

	// Set session cookie or token (example)
	http.SetCookie(w, &http.Cookie{
		Name:  "session",
		Value: sessionID,
		Path:  "/",
	})

	// Return success response
	responseData := map[string]string{"message": "Login successful"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseData)
}
