package controllers

import (
	"fmt"
	"main/database"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CompanyPage(c *gin.Context) {

	fmt.Println("Rendering company page")
	// Get logged user dat from context (set by JWT middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login?error=Usuário não autenticado")
		return
	}
	fmt.Println("Logged user ID:", userID)

	// You can use userID to fetch user details from the database if needed
	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/login?error=Usuário não autenticado")
		return
	}

	// retrieve company information if needed
	var company models.Company
	if err := database.DB.Where("id = ?", user.CompanyID).First(&company).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching company: %v", err)
		return
	}

	// Fetch user details from the database using userID if needed
	c.HTML(http.StatusOK, "company.html", gin.H{
		"title":        "Administre seu negócio!",
		"user_name":    user.Name,
		"user_email":   user.Email,
		"company_name": company.Name,
	})
}
