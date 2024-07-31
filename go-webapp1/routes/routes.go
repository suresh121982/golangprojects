package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"examples.com/webapp/cache"
	"examples.com/webapp/data"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

// Secret key for JWT signing
var jwtKey = []byte("your_secret_key") // Change this to a more secure key

// UserCredentials struct to hold login information
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Claims struct for JWT
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/login", Login).Methods("POST")
	r.HandleFunc("/items", GetItems).Methods("GET")
	r.HandleFunc("/items", CreateItem).Methods("POST")
}

func Login(w http.ResponseWriter, r *http.Request) {
	var creds UserCredentials
	fmt.Println("Received credentials:", creds)
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Log the received credentials
	fmt.Println("Received credentials:", creds)

	// Validate credentials (this is a placeholder; implement your logic)
	if creds.Username != "test" || creds.Password != "password" {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}

func GetItems(w http.ResponseWriter, r *http.Request) {
	if items, found := cache.Get("items"); found {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(items)
		return
	}
	w.Write([]byte("No items found"))
}

func CreateItem(w http.ResponseWriter, r *http.Request) {
	var item data.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// WaitGroup to manage Goroutines
	var wg sync.WaitGroup

	// Save item to cache
	wg.Add(1)
	go func() {
		defer wg.Done()
		existingItems, found := cache.Get("items")
		if found {
			currentItems := existingItems.([]data.Item)
			currentItems = append(currentItems, item)
			cache.Set("items", currentItems)
		} else {
			cache.Set("items", []data.Item{item})
		}
	}()

	// Save item to file
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := data.SaveToFile([]data.Item{item}, "items.json")
		if err != nil {
			// Handle error (you may want to log this)
			fmt.Println("Error saving to file:", err)
		}
	}()

	// Wait for both operations to complete
	wg.Wait()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
}
