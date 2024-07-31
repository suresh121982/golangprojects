// api.go
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"examples.com/go-webapp/internal/data"
)

// GetAllUsersHandler returns all users as JSON
func GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := data.GetAllUsers()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetUserByIDHandler retrieves a user by ID
func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := data.GetUserByID(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("User not found: %v", err), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// AddUserHandler adds a new user
func AddUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser data.User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	data.AddUser(newUser)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User added successfully"))
}

// UpdateUserHandler updates an existing user
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	var updatedUser data.User
	err = json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = data.UpdateUser(userID, updatedUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("User not found: %v", err), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated successfully"))
}

// DeleteUserHandler deletes a user by ID
func DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = data.DeleteUser(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("User not found: %v", err), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User deleted successfully"))
}
