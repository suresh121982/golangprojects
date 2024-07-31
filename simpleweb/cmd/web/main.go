package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

// Example initial data
var initialData = map[string]string{
	"message": "Hello, World!",
}

func main() {
	InitUsers()

	// Define HTTP routes
	http.HandleFunc("/login", LoginHandler)

	// Initialize data directory if it doesn't exist
	err := os.MkdirAll("internal/data", 0755)
	if err != nil {
		log.Fatalf("Error creating data directory: %v", err)
	}

	// Write initial data to file if file doesn't exist
	if _, err := os.Stat("internal/data/data.json"); os.IsNotExist(err) {
		if err := WriteDataToFile(initialData, "internal/data/data.json"); err != nil {
			log.Fatalf("Error creating initial data file: %v", err)
		}
		log.Println("Initial data file created")
	}

	// Define HTTP routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Example handler
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(initialData)
	})

	// Start HTTP server
	port := ":8080"
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// WriteDataToFile writes data to a JSON file
func WriteDataToFile(data interface{}, filename string) error {
	jsonData, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
