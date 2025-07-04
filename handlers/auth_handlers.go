package handlers

import (
	"net/http"
	"strings"
	"time"

	"example.com/m/v2/auth"
	"example.com/m/v2/database"
	"example.com/m/v2/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterUserCompanyRequest defines the expected JSON structure for registration
type RegisterUserCompanyRequest struct {
	Name                 string `json:"name" binding:"required"`
	Email                string `json:"email" binding:"required,email"`
	Password             string `json:"password" binding:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" binding:"required,eqfield=Password"`
	CompanyName          string `json:"company_name" binding:"required"`
}

// RegisterUserCompanyHandler handles new user and company registration.
// New users are created with "Manager" privileges and companies at level 0.
func RegisterUserCompanyHandler(c *gin.Context) {
	var req RegisterUserCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if email already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already registered"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking email: " + err.Error()})
		return
	}

	// Check if company name already exists
	var existingCompany models.Company
	if err := database.DB.Where("name = ?", req.CompanyName).First(&existingCompany).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Company name already registered"})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error checking company name: " + err.Error()})
		return
	}

	// Get "Manager" privilege
	var managerPrivilege models.Privilege
	if err := database.DB.Where("name = ?", models.PrivilegeManager).First(&managerPrivilege).Error; err != nil {
		// Privileges should be seeded at startup. If not found here, it's an issue.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Manager privilege not found. System may not be initialized correctly."})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create company
	company := models.Company{
		Name:         req.CompanyName,
		Level:        0, // Default level 0
		Active:       true,
		CreationDate: time.Now(), // Hooks will also set this
		LastModified: time.Now(), // Hooks will also set this
		// Other company fields like ZIPCode, Street, etc., can be added here or updated later
	}

	// Create user
	user := models.User{
		Name:         req.Name,
		Email:        req.Email,
		Password:     hashedPassword,
		PrivilegeID:  managerPrivilege.ID,
		CreationDate: time.Now(), // Hooks will also set this
		LastModified: time.Now(), // Hooks will also set this
	}

	// Use a transaction to ensure both company and user are created, or neither.
	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	if err := tx.Create(&company).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create company: " + err.Error()})
		return
	}

	user.CompanyID = company.ID
	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	// Generate JWT for the new user
	token, err := auth.GenerateJWT(&user, managerPrivilege.Name)
	if err != nil {
		// Log this error but proceed to return user info, as registration was successful
		// The user can try to log in manually if token generation fails
		c.JSON(http.StatusCreated, gin.H{
			"message": "User and company registered successfully. Token generation failed.",
			"user_id": user.ID,
			"email":   user.Email,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User and company registered successfully",
		"token":   token,
		"user": gin.H{
			"id":           user.ID,
			"name":         user.Name,
			"email":        user.Email,
			"privilege_id": user.PrivilegeID,
			"company_id":   user.CompanyID,
		},
		"company": gin.H{
			"id":   company.ID,
			"name": company.Name,
		},
	})
}

// LoginRequest defines the expected JSON structure for login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginHandler handles user login.
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user models.User
	// Preload Privilege to get Privilege.Name for JWT
	if err := database.DB.Preload("Privilege").Where("email = ?", req.Email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error: " + err.Error()})
		return
	}

	if !auth.CheckPasswordHash(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Ensure Privilege is loaded and Privilege.Name is available
	if user.Privilege.Name == "" {
		// This case should ideally not happen if Preload worked and data is consistent
		// Fetch privilege explicitly if it's missing
		var privilege models.Privilege
		if err := database.DB.First(&privilege, user.PrivilegeID).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not determine user privilege"})
			return
		}
		user.Privilege = privilege // Assign loaded privilege
	}

	token, err := auth.GenerateJWT(&user, user.Privilege.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":             user.ID,
			"name":           user.Name,
			"email":          user.Email,
			"privilege_id":   user.PrivilegeID,
			"privilege_name": user.Privilege.Name,
			"company_id":     user.CompanyID,
		},
	})
}
