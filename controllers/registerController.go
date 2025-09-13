package controllers

import (
	"fmt"
	"main/database"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterCompanyUser(c *gin.Context) {

	// Placeholder for company registration logic
	company := models.Company{}
	if err := c.ShouldBind(&company); err != nil {
		c.String(http.StatusBadRequest, "Error binding company data: %v", err)
		return
	}
	// Bind "company_zip" from form as "zip" to company
	company.CEP = c.PostForm("company_zip")
	company.Name = c.PostForm("company_name")
	company.Street = c.PostForm("company_street")
	company.Number = c.PostForm("company_number")
	company.Neighborhood = c.PostForm("company_neighborhood")
	company.City = c.PostForm("company_city")
	company.State = c.PostForm("company_state")
	company.Phone = c.PostForm("company_phone")
	// Save the company to the database
	if err := database.DB.Create(&company).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error saving company: %v", err)
		return
	}

	// Placeholder for user registration logic
	user := models.User{}
	if err := c.ShouldBind(&user); err != nil {
		c.String(http.StatusBadRequest, "Error binding user data: %v", err)
		return
	}
	user.Company = company // Associate the user with the company
	// Bind the "profile" field from the form to the user.Profile field
	privilege := c.Query("privilege")
	if privilege == "" {
		privilege = c.PostForm("privilege")
	}
	if privilege == "" {
		privilege = "user" // Default privilege if none provided
	}
	user.Name = c.PostForm("name")
	user.Email = c.PostForm("email")
	user.Password = c.PostForm("password")
	// For debugging purposes, print the privilege being assigned
	fmt.Println("Registering user with privilege:", privilege)
	var priv models.Privilege
	if err := database.DB.Where("slug = ?", privilege).First(&priv).Error; err != nil {
		c.String(http.StatusBadRequest, "Invalid privilege: %v", err)
		return
	}
	user.Privilege = priv
	// Save the user to the database
	if err := SaveUser(user); err != nil {
		c.String(http.StatusInternalServerError, "Error saving user: %v", err)
		return
	}

	// Redirect to welcome page after registration
	c.Redirect(http.StatusSeeOther, "/")
}

func RegisterPage(c *gin.Context) {
	fmt.Println("Acessou a página de registro")
	c.HTML(http.StatusOK, "register.html",
		gin.H{
			"Title": "Cadastro de Usuário",
		})
}
