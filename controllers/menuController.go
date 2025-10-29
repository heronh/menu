package controllers

import (
	"fmt"
	"main/database"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ViewCompanyMenu(c *gin.Context) {

	// Get all dishes from this company
	CompanyID, _ := c.Get("company_id")
	var dishes []models.Dish
	if err := database.DB.Preload("Images").Where("company_id = ?", CompanyID).Find(&dishes).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dishes"})
		return
	}

	for _, dish := range dishes {
		fmt.Println("\n", dish.Name)
		fmt.Println("Active:", dish.Active)
		fmt.Println(dish.Description)
		fmt.Println(dish.Price)
		fmt.Println("Images:")
		// You can access dish.Images here
		_ = dish.Images // Just to avoid unused variable warning
		for _, img := range dish.Images {
			fmt.Println(" - Image URL:", img.OriginalFileName)
		}
	}

	// Render the menu template
	c.HTML(http.StatusOK, "menu-company.html", gin.H{
		"Title":     "Menu",
		"Dishes":    dishes,
		"CompanyID": CompanyID,
	})
}

func ViewCustomerMenu(c *gin.Context) {

	CompanyID, _ := c.Get("company_id")

	// Get all dishes from this company
	var dishes []models.Dish
	if err := database.DB.Preload("Images").Where("company_id = ?", CompanyID).Find(&dishes).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dishes"})
		return
	}

	// Render the menu template
	c.HTML(http.StatusOK, "menu.html", gin.H{
		"Title":     "Menu",
		"Dishes":    dishes,
		"CompanyID": CompanyID,
	})
}
