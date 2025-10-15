package controllers

import (
	"bytes"
	"fmt"
	"main/database"
	"main/models"
	"net/http"
	"text/template"

	"github.com/gin-gonic/gin"
)

func DeleteSection(c *gin.Context) {
	sectionID := c.Param("id")
	if sectionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Section ID is required"})
		return
	}

	companyID, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in token"})
		return
	}

	var section models.Section
	if err := database.DB.Where("id = ? AND company_id = ?", sectionID, companyID.(uint)).First(&section).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	if err := database.DB.Delete(&section).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete section"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Section deleted successfully"})
}

func CreateSection(c *gin.Context) {
	type SectionInput struct {
		Description string `json:"description" binding:"required"`
	}
	fmt.Println("CreateSection called")

	var input SectionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Description is required"})
		return
	}
	fmt.Println("Input description:", input.Description)

	if len(input.Description) < 3 || len(input.Description) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Description must be between 3 and 100 characters"})
		return
	}

	if len(input.Description) > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Description must be less than 100 characters"})
		return
	}

	for _, char := range input.Description {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && (char < '0' || char > '9') && char != ' ' && char != '-' && char != '_' {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Description contains invalid characters"})
			return
		}
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}
	fmt.Println("User ID from token:", userId)

	companyID, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in token"})
		return
	}
	fmt.Println("Company ID from token:", companyID)

	section := models.Section{
		Description: input.Description,
		CompanyID:   companyID.(uint),
		AuthorID:    userId.(uint),
	}

	var existingSection models.Section
	if err := database.DB.Where("description = ? AND company_id = ?", input.Description, companyID.(uint)).First(&existingSection).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Esta seção já existe"})
		return
	}
	fmt.Println("No existing section found with the same description")

	if err := database.DB.Create(&section).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create section"})
		return
	}
	fmt.Println("Section created successfully with ID:", section.ID)

	renderedHTML, err := renderSectionBox(section)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to render section HTML"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Section created successfully", "renderedHTML": renderedHTML})
}

func renderSectionBox(section models.Section) (string, error) {

	fmt.Println("Rendering section box for section:", section.Description)
	fmt.Println("Section ID:", section.ID)
	fmt.Println("Section Description:", section.Description)

	var tmpl, err = template.ParseFiles("templates/dish/box-section.html")
	if err != nil {
		fmt.Println("Error loading template:", err)
		return "", err
	}
	fmt.Println("Template loaded successfully")

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, section); err != nil {
		fmt.Println("Error rendering template:", err)
		return "", err
	}
	return buf.String(), nil
}
