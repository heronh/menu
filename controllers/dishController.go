package controllers

import (
	"fmt"
	"main/database"
	"main/models"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

var TimeList = []string{
	"00:00", "00:30",
	"01:00", "01:30",
	"02:00", "02:30",
	"03:00", "03:30",
	"04:00", "04:30",
	"05:00", "05:30",
	"06:00", "06:30",
	"07:00", "07:30",
	"08:00", "08:30",
	"09:00", "09:30",
	"10:00", "10:30",
	"11:00", "11:30",
	"12:00", "12:30",
	"13:00", "13:30",
	"14:00", "14:30",
	"15:00", "15:30",
	"16:00", "16:30",
	"17:00", "17:30",
	"18:00", "18:30",
	"19:00", "19:30",
	"20:00", "20:30",
	"21:00", "21:30",
	"22:00", "22:30",
	"23:00", "23:30",
}

func NewDishPage(c *gin.Context) {

	//userID, _ := c.Get("user_id")
	companyID, _ := c.Get("company_id")
	userID, _ := c.Get("user_id")

	Sections := []models.Section{}
	if err := database.DB.Where("company_id = ?", companyID).Find(&Sections).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching sections: %v", err)
		return
	}
	fmt.Println("Total sections found:", len(Sections))

	var Images []models.Image
	if err := database.DB.Where("company_id = ?", companyID).Find(&Images).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching images: %v", err)
		return
	}
	fmt.Println("Total images found:", len(Images))

	for _, img := range Images {
		fmt.Println("Image:", img.ID, img.OriginalFileName, img.UniqueName)
		// Check if file exists
		if _, err := os.Stat(img.UniqueName); os.IsNotExist(err) {
			fmt.Println("File does not exist:", img.UniqueName)
		} else {
			fmt.Println("File exists:", img.UniqueName)
		}
	}

	c.HTML(http.StatusOK, "new-dish.html", gin.H{
		"title":     "Adicionar novo prato",
		"Sections":  Sections,
		"Images":    Images,
		"CompanyId": companyID,
		"UserId":    userID,
		"TimeList":  TimeList,
	})
}

func CreateDish(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login?error=Usuário não autenticado")
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/login?error=Usuário não autenticado")
		return
	}

	var company models.Company
	if err := database.DB.Where("id = ?", user.CompanyID).First(&company).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching company: %v", err)
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Preço inválido: %v", err)
		return
	}

	dish := models.Dish{
		Name:        name,
		Description: description,
		Price:       price,
		CompanyID:   company.ID,
	}

	if err := database.DB.Create(&dish).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error creating dish: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/company")
}

func EditDishPage(c *gin.Context) {
	id := c.Param("id")
	var dish models.Dish
	if err := database.DB.Where("id = ?", id).First(&dish).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching dish: %v", err)
		return
	}

	var Images []models.Image
	if err := database.DB.Where("company_id = ?", dish.CompanyID).Find(&Images).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching images: %v", err)
		return
	}
	for _, img := range Images {
		fmt.Println("Image:", img.ID, img.OriginalFileName, img.UniqueName)
		// Check if file exists
		if _, err := os.Stat(img.UniqueName); os.IsNotExist(err) {
			fmt.Println("File does not exist:", img.UniqueName)
		} else {
			fmt.Println("File exists:", img.UniqueName)
		}
	}

	c.HTML(http.StatusOK, "edit_dish.html", gin.H{
		"title":  "Editar prato",
		"dish":   dish,
		"Images": Images,
	})
}

func UpdateDish(c *gin.Context) {
	id := c.Param("id")
	var dish models.Dish
	if err := database.DB.Where("id = ?", id).First(&dish).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching dish: %v", err)
		return
	}

	name := c.PostForm("name")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Preço inválido: %v", err)
		return
	}

	dish.Name = name
	dish.Description = description
	dish.Price = price

	if err := database.DB.Save(&dish).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error updating dish: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/company")
}

func DeleteDish(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.Dish{}, id).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error deleting dish: %v", err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/company")
}

func ValidateDish(c *gin.Context) {
	// Example validation: check if name is provided
	name := c.PostForm("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"valid": false, "error": "Name is required"})
		return
	}
	// Additional validations can be added here

	c.JSON(http.StatusOK, gin.H{"valid": true})
}
