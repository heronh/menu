package controllers

import (
	"main/database"
	"main/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewDishPage(c *gin.Context) {

	sections := []models.Section{}
	database.DB.Find(&sections)

	// Popular sections para testes
	sections = append(sections, models.Section{ID: 1, Description: "Entradas"})
	sections = append(sections, models.Section{ID: 2, Description: "Pratos Principais"})
	sections = append(sections, models.Section{ID: 3, Description: "Sobremesas"})
	sections = append(sections, models.Section{ID: 4, Description: "Bebidas"})

	c.HTML(http.StatusOK, "new-dish.html", gin.H{
		"title":    "Adicionar novo prato",
		"Sections": sections,
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
	c.HTML(http.StatusOK, "edit_dish.html", gin.H{
		"title": "Editar prato",
		"dish":  dish,
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
