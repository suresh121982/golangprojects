package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// User struct represents a user in the system.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Credentials struct represents the login credentials received from clients.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var jwtKey = []byte("your-secret-key") // Change this to your secret key

// UserCache to store logged in users
var userCache = make(map[string]time.Time)

// Check if user is authenticated
func isAuthenticated(username string) bool {
	_, ok := userCache[username]
	return ok
}

// Login handler generates a JWT for a user
func login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate credentials (for simplicity, using a hardcoded check)
	validUser := User{
		Username: "testuser",
		Password: "password123",
	}

	if creds.Username != validUser.Username || creds.Password != validUser.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create JWT token
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   creds.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Store user in cache
	userCache[creds.Username] = time.Now()
	saveUserCacheToFile()
	// Append user session to sessions.json
	appendSessionToFile(creds.Username)

	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	})

	w.WriteHeader(http.StatusOK)
}

// AccessRestricted handler to access restricted endpoint
func accessRestricted(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := cookie.Value
	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Check if user is authenticated
	if !isAuthenticated(claims.Subject) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.Write([]byte(fmt.Sprintf("Welcome, %s!", claims.Subject)))
}

// Logout handler to logout user and remove from cache
func logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	tokenString := cookie.Value
	claims := &jwt.StandardClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Remove user from cache
	delete(userCache, claims.Subject)
	saveUserCacheToFile()

	// Clear cookie
	http.SetCookie(w, &http.Cookie{
		Name:   "jwt",
		Value:  "",
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusOK)
}

func adminResetSessionsHandler(w http.ResponseWriter, r *http.Request) {
	// Check if user is authorized (e.g., admin role)
	// Example: Replace with your own authorization logic
	if !isAdmin(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Reset sessions
	if err := resetSessions(); err != nil {
		http.Error(w, "Error resetting sessions", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	fmt.Fprintf(w, "Sessions reset successfully.")
}

func resetSessions() error {
	// Path to the sessions file
	sessionsFile := "sessions.json"

	// Check if the sessions file exists
	_, err := os.Stat(sessionsFile)
	if os.IsNotExist(err) {
		// If the file doesn't exist, return nil (no error)
		return nil
	} else if err != nil {
		// If there's an error other than file not exists, return the error
		return fmt.Errorf("error checking sessions file: %v", err)
	}

	// Attempt to remove the sessions file
	err = os.Remove(sessionsFile)
	if err != nil {
		return fmt.Errorf("error deleting sessions file: %v", err)
	}

	// Print a message indicating success
	fmt.Println("Sessions reset successfully.")

	return nil
}

func ensureSessionsFileExists(sessionsFile string) error {
	// Check if the sessions file exists
	_, err := os.Stat(sessionsFile)
	if os.IsNotExist(err) {
		// If the file doesn't exist, create it
		file, err := os.Create(sessionsFile)
		if err != nil {
			return fmt.Errorf("error creating sessions file: %v", err)
		}
		defer file.Close() // Ensure file is closed after function execution
	} else if err != nil {
		// If there's an error other than file not exists, return the error
		return fmt.Errorf("error checking sessions file: %v", err)
	}

	return nil
}

func isAdmin(r *http.Request) bool {
	// Example function to check if user is an admin based on some criteria
	// Replace with your own implementation
	// Example: check for admin role in JWT claims, database lookup, etc.
	return true // Replace with actual admin check logic
}

func main() {
	router := mux.NewRouter()
	loadUserCacheFromFile()

	sessionsFile := "sessions.json"

	// Ensure sessions file exists (or is created if it doesn't exist)
	if err := ensureSessionsFileExists(sessionsFile); err != nil {
		fmt.Printf("Failed to ensure sessions file exists: %v\n", err)
		return
	}

	// Now you can proceed with other operations involving sessionsFile
	fmt.Println("Sessions file is ready for use.")

	// Route to handle user login
	router.HandleFunc("/login", login).Methods("POST")

	// Route to handle restricted access
	router.HandleFunc("/restricted", accessRestricted).Methods("GET")

	// Route to handle user logout
	router.HandleFunc("/logout", logout).Methods("POST")

	router.HandleFunc("/admin/reset-sessions", adminResetSessionsHandler).Methods("POST")

	// Start server
	fmt.Println("Server running on port 8000")
	err := http.ListenAndServe(":8000", router)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Load user sessions from file
func loadUserCacheFromFile() {
	file, err := ioutil.ReadFile("sessions.json")
	if err != nil {
		fmt.Println("Error reading sessions file:", err)
		return
	}

	err = json.Unmarshal(file, &userCache)
	if err != nil {
		fmt.Println("Error unmarshalling sessions:", err)
	}
}

// Save user sessions to file
// Example function to save userCache to sessions.json
// Function to save user sessions to file
func saveUserCacheToFile() {
	file, err := os.OpenFile("sessions.json", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("Error opening sessions file:", err)
		return
	}
	defer file.Close()

	// Convert userCache to JSON and write to file
	data, err := json.MarshalIndent(userCache, "", "  ")
	if err != nil {
		fmt.Println("Error marshalling userCache:", err)
		return
	}

	_, err = file.Write(data)
	if err != nil {
		fmt.Println("Error writing to sessions file:", err)
	}
}

func appendSessionToFile(username string) {
	file, err := os.OpenFile("sessions.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening sessions file:", err)
		return
	}
	defer file.Close()

	// Append username to file
	if _, err := file.WriteString(username + "\n"); err != nil {
		fmt.Println("Error appending to sessions file:", err)
	}
}

func deleteSessionFile() error {
	err := os.Remove("sessions.json")
	if err != nil {
		return err
	}
	return nil
}

// Modify login, logout, and isAuthenticated functions to read from/write to file
