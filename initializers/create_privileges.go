package initializers

import (
	"log"
	"main/database"
	"main/models"

	"gorm.io/gorm"
)

func CreatePrivileges() {
	if database.DB == nil {
		log.Fatal("Database connection not initialized before CreatePrivileges")
		return
	}

	privileges := []models.Privilege{
		{Name: "Super Administrator", Slug: "su"},
		{Name: "Administrator", Slug: "admin"},
		{Name: "Manager", Slug: "manager"},
		{Name: "Employee", Slug: "employee"},
	}

	for _, p := range privileges {
		var existingPrivilege models.Privilege
		// Check if privilege already exists by slug
		result := database.DB.Where("slug = ?", p.Slug).First(&existingPrivilege)
		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				// Privilege does not exist, create it
				if err := database.DB.Create(&p).Error; err != nil {
					log.Fatalf("Failed to create privilege %s: %v", p.Name, err)
				}
				log.Printf("Privilege %s created successfully", p.Name)
			} else {
				// Some other error occurred
				log.Fatalf("Failed to query privilege %s: %v", p.Name, result.Error)
			}
		} else {
			log.Printf("Privilege %s already exists", p.Name)
		}
	}
	log.Println("Privilege creation process completed.")
}
