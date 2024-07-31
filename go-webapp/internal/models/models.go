// user.go

package models

import (
	"errors"
	// Import necessary packages
)

// User model
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	// Add other fields as needed
}

// Mock user data for example
var users = []User{
	{ID: 1, Username: "user1", Role: "admin"},
	{ID: 2, Username: "user2", Role: "user"},
}

// GetUserByID retrieves a user by ID
func GetUserByID(userID int) (*User, error) {
	for _, user := range users {
		if user.ID == userID {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}
