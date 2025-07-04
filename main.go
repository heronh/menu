package main

import (
	"log"

	"example.com/m/v2/auth"
	"example.com/m/v2/database"
	"example.com/m/v2/routes"
	"example.com/m/v2/seed" // Import the seed package
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize authentication module (e.g., load JWT secret key)
	auth.Initialize()
	log.Println("Authentication module initialized.")

	// Initialize database connection
	database.Connect() // This also runs AutoMigrate
	log.Println("Successfully connected to and migrated the database.")

	// Seed initial data (Privileges, Super Admin User)
	// This should run after migrations to ensure tables exist.
	seed.SeedData()

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
