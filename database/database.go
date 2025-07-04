package database

import (
	"fmt"
	"log"
	"os"

	"example.com/m/v2/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Connect initializes the database connection and performs auto-migration.
func Connect() {
	var err error
	// Connection string - replace with your actual connection details
	// It's good practice to use environment variables for sensitive data
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=gin_gorm_app port=5432 sslmode=disable TimeZone=Asia/Shanghai"
		log.Println("DATABASE_URL environment variable not set, using default DSN.")
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	fmt.Println("Database connection successfully opened")

	// Auto-migrate the schema
	// This will ONLY create tables, missing columns and missing indexes.
	// It will NOT change existing column types or delete unused columns to protect your data.
	err = DB.AutoMigrate(
		// Order matters if there are strict foreign key dependencies not automatically handled,
		// but GORM is generally good at figuring out the order.
		// Migrating models with fewer dependencies first can sometimes be safer.
		&models.Privilege{},
		&models.Company{},  // Company might be referenced by User, Dish, Image, Message
		&models.User{},     // User references Privilege and Company. User is referenced by many other models.
		&models.Category{}, // Category is referenced by Dish
		&models.Image{},    // Image references User and Company
		&models.Dish{},     // Dish references Category, User (author), Company
		&models.Message{},  // Message references User (sender, recipient) and Company (sender company)
		&models.Log{},
		&models.Todo{}, // Todo references User
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	fmt.Println("Database migrated")
}
