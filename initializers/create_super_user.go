package initializers

import (
	"log"
	"main/database"
	"main/models"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateSuperUser() {
	if database.DB == nil {
		log.Fatal("Database connection not initialized before CreateSuperUser")
		return
	}

	suEmail := os.Getenv("USER_SEEDER_EMAIL")
	suPassword := os.Getenv("USER_SEEDER_PASSWORD")
	suName := os.Getenv("USER_SEEDER_NAME")

	if suEmail == "" || suPassword == "" || suName == "" {
		log.Println("USER_SEEDER_EMAIL, USER_SEEDER_PASSWORD, or USER_SEEDER_NAME environment variables not set. Skipping super user creation.")
		return
	}

	// Check if super user already exists
	var existingUser models.User
	err := database.DB.Where("email = ?", suEmail).First(&existingUser).Error
	if err == nil {
		log.Printf("Super user with email %s already exists. Skipping creation.", suEmail)
		return
	}
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Fatalf("Failed to query for existing super user: %v", err)
		return
	}

	// Find the Super Administrator privilege
	var suPrivilege models.Privilege
	if err := database.DB.Where("slug = ?", "su").First(&suPrivilege).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			log.Fatal("Super Administrator privilege not found. Please run CreatePrivileges first.")
		} else {
			log.Fatalf("Failed to find Super Administrator privilege: %v", err)
		}
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(suPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Failed to hash super user password: %v", err)
		return
	}

	superUser := models.User{
		Name:        suName,
		Email:       suEmail,
		Password:    string(hashedPassword),
		PrivilegeID: suPrivilege.ID,
		CompanyID:   nil, // Null company
	}

	if err := database.DB.Create(&superUser).Error; err != nil {
		log.Fatalf("Failed to create super user: %v", err)
	}

	log.Printf("Super user %s <%s> created successfully.", suName, suEmail)
}
