// main.go
package main

import (
	"log"
	"net/http"

	"examples.com/go-webapp/internal/data"
	"examples.com/go-webapp/internal/handlers"
)

func main() {
	// Initialize sample users (for demo purposes)
	data.InitUsers()

	// Define API endpoints
	http.HandleFunc("/api/users", handlers.GetAllUsersHandler)
	http.HandleFunc("/api/user", handlers.GetUserByIDHandler)
	http.HandleFunc("/api/adduser", handlers.AddUserHandler)
	http.HandleFunc("/api/updateuser", handlers.UpdateUserHandler)
	http.HandleFunc("/api/deleteuser", handlers.DeleteUserHandler)

	// Serve static files (e.g., CSS)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start HTTP server
	port := ":8080"
	log.Printf("Server listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
