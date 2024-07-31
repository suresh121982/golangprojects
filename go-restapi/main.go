// main.go
package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"examples.com/go-restapi/handlers"
	"examples.com/go-restapi/middleware"
)

func main() {
	// Load configurations
	err := handlers.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize Gin router
	router := gin.Default()

	// Setup middleware
	router.Use(middleware.CORSMiddleware())

	// Setup routes
	handlers.SetupAuthRoutes(router)
	handlers.SetupItemRoutes(router)

	// Start server
	port := viper.GetString("server.port")
	log.Printf("Starting server on port %s", port)
	err = router.Run(":" + port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
