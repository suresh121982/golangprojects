package main

import (
	"examples.com/webapp/middleware"
	"examples.com/webapp/routes"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.Use(middleware.Auth) // Use authorization middleware
	routes.RegisterRoutes(r)

	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
