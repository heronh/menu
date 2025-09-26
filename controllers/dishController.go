package controllers

import (
	"fmt"
	"main/database"
	"main/models"
	"net/http"
	"os"
	"strconv"
	"time"

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

func UploadDishImage(c *gin.Context) {
	dishID := c.Param("id")

	// Fetch the dish to ensure it exists
	var dish models.Dish
	if err := database.DB.Where("id = ?", dishID).First(&dish).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching dish: %v", err)
		return
	}

	// Retrieve the file from the form data
	file, err := c.FormFile("image")
	if err != nil {
		c.String(http.StatusBadRequest, "Error retrieving file: %v", err)
		return
	}

	// Save the file to a specific location (e.g., "./uploads/")
	filePath := "./uploads/" + file.Filename
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.String(http.StatusInternalServerError, "Error saving file: %v", err)
		return
	}
	// Create a new Image record
	dishImage := models.Image{
		ID:               dish.ID,
		OriginalFileName: filePath, // In a real app, this would be a URL accessible by the frontend
	}
	if err := database.DB.Create(&dishImage).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error saving image record: %v", err)
		return
	}
	c.Redirect(http.StatusSeeOther, "/dishes/edit/"+dishID)
}

func DeleteDishImage(c *gin.Context) {
	dishID := c.Param("dish_id")
	imageID := c.Param("image_id")

	// Fetch the dish image to ensure it exists
	var dishImage models.Image
	if err := database.DB.Where("id = ? AND dish_id = ?", imageID, dishID).First(&dishImage).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching dish image: %v", err)
		return
	}

	// Delete the dish image record
	if err := database.DB.Delete(&dishImage).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error deleting dish image: %v", err)
		return
	}

	c.Redirect(http.StatusSeeOther, "/dishes/edit/"+dishID)
}

func UploadMultipleDishImages(c *gin.Context) {

	// list of files
	fmt.Println("Uploading multiple images...")

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("Erro ao ler o formulário: %s", err.Error()),
		})
		c.Redirect(http.StatusFound, c.Request.Referer())
		return
	}

	// Create a subfolder based on current user or company
	userID, _ := c.Get("user_id")
	companyID, _ := c.Get("company_id")
	uploadPath := fmt.Sprintf("./uploads/company_%d/user_%d/", companyID, userID)
	os.MkdirAll(uploadPath, os.ModePerm)

	files := form.File["images[]"]
	for _, file := range files {
		// Save the file to a specific location (e.g., "./uploads/")
		fmt.Println("Uploading file:", file.Filename)
		//Limit file types to images only
		if file.Header.Get("Content-Type") != "image/jpeg" && file.Header.Get("Content-Type") != "image/png" && file.Header.Get("Content-Type") != "image/gif" {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": "Apenas arquivos de imagem (JPEG, PNG, GIF) são permitidos.",
			})
			return
		}

		// Save the file
		fmt.Println("Saving file to:", uploadPath+file.Filename)
		if err := c.SaveUploadedFile(file, uploadPath+file.Filename); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": fmt.Sprintf("Erro ao salvar o arquivo: %s", err.Error()),
			})
			return
		}

		// Update image database
		fmt.Println("Creating image record in database...")
		fmt.Println("For user ID:", userID)
		fmt.Println("For company ID:", companyID)
		dishImage := models.Image{
			OriginalFileName: file.Filename,
			UniqueName:       uploadPath + file.Filename,
			Storage:          "local",
			InsertedByID:     userID.(uint),
			CompanyID:        companyID.(uint),
			CreatedAt:        time.Now(),
		}
		fmt.Println("Saving image record to database:", dishImage)
		fmt.Println("With path:", uploadPath+file.Filename)
		fmt.Println("At time:", time.Now())
		// Save record to database
		if err := database.DB.Create(&dishImage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": fmt.Sprintf("Erro ao salvar o registro da imagem: %s", err.Error()),
			})
			return
		}
	}
	c.Redirect(http.StatusFound, c.Request.Referer())
}

/*
type Image struct {
	gorm.Model
	ID               uint    `json:"id" gorm:"primary_key"`
	OriginalFileName string  `gorm:"not null"`
	UniqueName       string  `gorm:"unique;not null"`
	Storage          string  // e.g., local, s3
	InsertedByID     uint    `gorm:"not null"`
	InsertedBy       User    `gorm:"foreignKey:InsertedByID"`
	CompanyID        uint    `gorm:"not null"`
	Company          Company `gorm:"foreignKey:CompanyID"`
	CreatedAt        time.Time
*/
