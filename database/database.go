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
		&models.User{},
		&models.Company{},
		&models.Role{},
		&models.DishCategory{},
		&models.Image{},
		&models.Dish{}, // Dish model depends on DishCategory and Image, so migrate them first or ensure GORM handles dependencies.
	)
	if err != nil {
		log.Fatalf("Failed to auto-migrate database: %v", err)
	}
	fmt.Println("Database migrated")
}
