package main

import (
	"log"
	"main/database"
	"main/initializers"
)

func init() {
	// Load .env file (if any)
	initializers.LoadEnvVariables()

	// Connect to DB
	database.Connect()

	// Run Migrations
	initializers.SyncDatabase()

	// Create Privileges
	initializers.CreatePrivileges()

	// Create Super User
	initializers.CreateSuperUser()
}

func main() {
	log.Println("Application setup complete. Starting application...")
	// Application logic will go here in the future.
	// For now, it just logs that setup is done.
	// Example: You might start an HTTP server here.
	// e.g., router := setupRouter()
	// router.Run(":8080")
}
