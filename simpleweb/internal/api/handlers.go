package api

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"

	"examples.com/myapp/internal/cache"
)

// HTML rendering function
func RenderIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("internal/views/index.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Example data
	data := map[string]string{"Title": "Welcome to MyApp", "Message": "Hello, World!"}

	// Render HTML with data
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// REST API function
func GetData(w http.ResponseWriter, r *http.Request) {
	// Fetch data logic
	data := map[string]string{"message": "Hello, World!"}

	// Write to cache and file
	cache.Set("data", data)
	writeToFile(data)

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func writeToFile(data map[string]string) {
	file, err := os.OpenFile("internal/data/data.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error writing to file:", err)
	} else {
		defer file.Close()
		if err := json.NewEncoder(file).Encode(data); err != nil {
			log.Println("Error encoding JSON to file:", err)
		}
	}
}

// Middleware function
