package auth

import (
	"encoding/json"
	"net/http"
	"time"

	"exmaples.com/myapp/logging"
	"github.com/dgrijalva/jwt-go"
)

// JWT secret key
var jwtKey = []byte("your-secret-key")

// Credentials struct represents the data received from client during login
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Function to authenticate user (replace with your actual authentication logic)
func authenticate(username, password string) bool {
	// Example: hardcoded authentication
	return username == "admin" && password == "admin"
}

// LoginHandler generates a JWT for a user
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		logging.Log("ERROR", "Failed to decode request body: "+err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !authenticate(creds.Username, creds.Password) {
		logging.Log("WARNING", "Failed authentication attempt for user: "+creds.Username)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   creds.Username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		logging.Log("ERROR", "Failed to sign JWT token: "+err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set JWT token as cookie in response
	http.SetCookie(w, &http.Cookie{
		Name:     "jwt",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
	})
	logging.Log("INFO", "User logged in successfully: "+creds.Username)
	w.WriteHeader(http.StatusOK)
}

// AccessRestrictedHandler handles access to restricted endpoints
func AccessRestrictedHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Access restricted\n"))
}

// Middleware function to check JWT token and authorize user
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			if err == http.ErrNoCookie {
				logging.Log("WARNING", "No JWT cookie found")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			logging.Log("ERROR", "Failed to get JWT cookie: "+err.Error())
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
				logging.Log("WARNING", "Invalid JWT signature")
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			logging.Log("ERROR", "Failed to parse JWT token: "+err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if !token.Valid {
			logging.Log("WARNING", "Invalid JWT token")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Check if user is authenticated
		// Replace with your actual logic to check if user exists in cache, database, etc.
		if !authenticate(claims.Subject, "") {
			logging.Log("WARNING", "Failed authentication check for user: "+claims.Subject)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		logging.Log("INFO", "User authorized successfully: "+claims.Subject)
		next.ServeHTTP(w, r)
	}
}
