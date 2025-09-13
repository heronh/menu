package controllers

import (
	"crypto/sha256"
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

func SaveUser(user models.User) error {

	// Hash the user's password before saving
	user.Password = HashPassword(user.Password)
	if err := database.DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
