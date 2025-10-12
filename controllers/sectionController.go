package controllers

import (
	"main/database"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateSection(c *gin.Context) {
	var input struct {
		Description string `json:"description" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	companyID, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Company ID not found in token"})
		return
	}

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
	if err := database.DB.Create(&section).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create section"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": "Section created successfully", "section": section})
}

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

	c.JSON(http.StatusOK, gin.H{"message": "Section deleted successfully"})
}
