package initializers

import (
	"log"
	"main/database"
	"main/models"
)

func SyncDatabase() {
	if database.DB == nil {
		log.Fatal("Database connection not initialized before SyncDatabase")
		return
	}
	log.Println("Running Migrations")
	err := database.DB.AutoMigrate(
		&models.Privilege{},
		&models.User{},
		&models.Company{},
		&models.Category{},
		&models.Dish{},
		&models.Image{},
		&models.Message{},
		&models.Log{},
		&models.Todo{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrated successfully")
}
