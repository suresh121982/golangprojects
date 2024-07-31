// user.go

package models

import "errors"

// User represents a user entity
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
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

// HasPermission checks if the user has the specified permission
func (u *User) HasPermission(requiredPermission string) bool {
	// Example logic: Check if the user's role matches the required permission
	return u.Role == "admin" || u.Role == "manager" || u.Role == "editor"
}
