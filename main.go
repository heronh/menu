package main

import (
	"log"

	"example.com/m/v2/database"
	"example.com/m/v2/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database connection
	database.Connect()
	log.Println("Successfully connected to and migrated the database.")

	// Set up Gin router
	router := routes.SetupRouter()
	log.Println("Gin router setup complete.")

	// Start the server
	port := "8080" // You can make this configurable, e.g., via environment variables
	log.Printf("Starting server on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
