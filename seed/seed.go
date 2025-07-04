package seed

import (
	"log"

	"example.com/m/v2/auth" // For password hashing
	"example.com/m/v2/database"
	"example.com/m/v2/models"
	"gorm.io/gorm"
)

// SeedData checks for and creates initial data like privileges and a super admin user.
func SeedData() {
	log.Println("Starting data seeding process...")

	// Seed Privileges
	seedPrivileges()

	// Seed Super Administrator User
	seedSuperAdminUser()

	log.Println("Data seeding process completed.")
}

func seedPrivileges() {
	privileges := []models.Privilege{
		{Name: models.PrivilegeSuperAdministrator},
		{Name: models.PrivilegeManager},
		{Name: models.PrivilegeEmployee},
	}

	for _, p := range privileges {
		var existingPrivilege models.Privilege
		err := database.DB.Where("name = ?", p.Name).First(&existingPrivilege).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if createErr := database.DB.Create(&p).Error; createErr != nil {
					log.Printf("Failed to seed privilege '%s': %v\n", p.Name, createErr)
				} else {
					log.Printf("Privilege '%s' seeded successfully.\n", p.Name)
				}
			} else {
				log.Printf("Error checking privilege '%s': %v\n", p.Name, err)
			}
		} else {
			log.Printf("Privilege '%s' already exists.\n", p.Name)
		}
	}
}

func seedSuperAdminUser() {
	suEmail := "su@example.com" // Default SU email, consider making this configurable

	// Check if SU user already exists
	var existingUser models.User
	err := database.DB.Where("email = ?", suEmail).First(&existingUser).Error
	if err == nil {
		log.Printf("Super admin user '%s' already exists.\n", suEmail)
		return // User exists, no need to seed
	}
	if err != gorm.ErrRecordNotFound {
		log.Printf("Error checking for super admin user '%s': %v\n", suEmail, err)
		return // Some other DB error
	}

	// SU user does not exist, proceed to create

	// Get Super Administrator privilege
	var suPrivilege models.Privilege
	if err := database.DB.Where("name = ?", models.PrivilegeSuperAdministrator).First(&suPrivilege).Error; err != nil {
		log.Printf("Failed to find Super Administrator privilege for seeding SU user: %v. Ensure privileges are seeded first.\n", err)
		return
	}

	// Create a default "System" or "Admin" company for the SU, or leave CompanyID as 0/nil if SU is not tied to a specific company.
	// For this setup, an SU might not belong to an operational company, or might belong to a special "System" company.
	// Let's create a placeholder "System Administration" company for the SU.
	var systemCompany models.Company
	companyName := "System Administration HQ"
	err = database.DB.Where("name = ?", companyName).First(&systemCompany).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			systemCompany = models.Company{
				Name:   companyName,
				Level:  100, // Max level for system company
				Active: true,
				CNPJ:   "SYSTEM_ADMIN_CNPJ", // Unique placeholder
				// Other fields can be defaults or empty
			}
			if createErr := database.DB.Create(&systemCompany).Error; createErr != nil {
				log.Printf("Failed to create system company for SU: %v\n", createErr)
				// Decide if SU creation should proceed without a company or fail.
				// For now, let's proceed without associating if company creation fails.
				// A better approach might be to ensure company creation or handle it.
			} else {
				log.Printf("Company '%s' seeded for SU.\n", companyName)
			}
		} else {
			log.Printf("Error checking for system company '%s': %v\n", companyName, err)
		}
	}

	// Default password for SU, should be changed immediately after first login.
	// Consider making this configurable via ENV var for first setup.
	suPassword := "superadmin123"
	hashedPassword, hashErr := auth.HashPassword(suPassword)
	if hashErr != nil {
		log.Printf("Failed to hash password for SU user: %v\n", hashErr)
		return
	}

	superAdmin := models.User{
		Name:        "Super Administrator",
		Email:       suEmail,
		Password:    hashedPassword,
		PrivilegeID: suPrivilege.ID,
		CompanyID:   systemCompany.ID, // Assign if company was created/found
		// CreationDate and LastModified will be set by BeforeCreate hook
	}

	if err := database.DB.Create(&superAdmin).Error; err != nil {
		log.Printf("Failed to seed super admin user '%s': %v\n", suEmail, err)
	} else {
		log.Printf("Super admin user '%s' with password '%s' seeded successfully. CompanyID: %d\n", suEmail, suPassword, systemCompany.ID)
		log.Println("IMPORTANT: Change the default Super Administrator password immediately after first login!")
	}
}
