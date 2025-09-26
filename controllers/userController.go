package controllers

import (
	"crypto/sha256"
	"fmt"
	"main/database"
	"main/models"
	"main/utils"
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

func LoginPage(c *gin.Context) {

	// For test purposes, list all users in the console
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		fmt.Println("Error fetching users:", err)
	}

	// Check for error message in query parameters
	errorMessage := c.Query("error")
	// Remove any existing error message from the URL after displaying it
	if errorMessage != "" {
		c.Request.URL.RawQuery = ""
	}

	// Get a random user from the database (for demonstration) witch privilege_id != 1
	var user models.User
	if err := database.DB.Where("privilege_id != ?", 1).Order("RANDOM()").First(&user).Error; err != nil {
		fmt.Println("Error fetching random user:", err)
	} else {
		fmt.Println("Random user fetched for demo:", user.Email)
	}

	// Render the login page with users and error message (if any)
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":    "Acesse sua conta",
		"users":    users,
		"error":    errorMessage,
		"email":    user.Email,
		"password": "123456", // Default password for demo purposes
	})
}

func LoginUser(c *gin.Context) {

	email := c.PostForm("email")
	password := c.PostForm("password")
	hashedPassword := HashPassword(password)

	var user models.User
	if err := database.DB.Where("email = ? AND password = ?", email, hashedPassword).First(&user).Error; err != nil {
		c.Redirect(http.StatusSeeOther, "/login?error=Credenciais inválidas, senha ou email incorretos")
		return
	}

	// Successful login: generate JWT with companyID
	companyID := uint(0)
	if user.CompanyID != nil {
		companyID = *user.CompanyID
	}
	token, err := utils.GenerateJWT(user.ID, companyID)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login?error=Erro ao gerar token JWT")
		return
	}
	c.SetCookie("token", token, 3600*24, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/company")
}

func LogoutUser(c *gin.Context) {
	// Clear the JWT cookie by setting its MaxAge to -1
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}
