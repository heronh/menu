package controllers

import (
	"fmt"
	"main/database"
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IsEmailAvailable(c *gin.Context) {

	email := c.Query("email")
	if email == "" {
		email = c.PostForm("email")
	}
	fmt.Println("Checking email availability for:", email)
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuário não encontrado", "exists": false})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuário encontrado", "exists": true})
}
