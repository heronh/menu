package controllers

import (
	"bytes"
	"fmt"
	"main/database"
	"main/models"
	"net/http"
	"os"
	"strconv"
	"text/template"
	"time"

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
	uploadPath := fmt.Sprintf("/uploads/_%d/", companyID)
	os.MkdirAll(uploadPath, os.ModePerm)
	fmt.Println("Upload path:", uploadPath)

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
		if err := c.SaveUploadedFile(file, "."+uploadPath+file.Filename); err != nil {
			fmt.Println("Error saving file:", err)
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

func ListDishImages(c *gin.Context) {

	fmt.Println("Listing images...")
	companyID, _ := c.Get("company_id")
	fmt.Println("For company ID:", companyID)
	var Images []models.Image
	if err := database.DB.Where("company_id = ?", companyID).Find(&Images).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Erro ao buscar imagens: %s", err.Error()),
		})
		return
	}

	var renderedImageBoxes string
	for _, img := range Images {
		fmt.Println("Rendering image box for image ID:", img.ID)
		renderedHTML, err := renderImageBox(img)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": fmt.Sprintf("Erro ao renderizar imagem: %s", err.Error()),
			})
			return
		}
		renderedImageBoxes += renderedHTML
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"html":    renderedImageBoxes,
	})
}

func CreateImageBox(c *gin.Context) {
	imageIDStr := c.PostForm("image_id")
	imageID, err := strconv.Atoi(imageIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID da imagem inválido.",
		})
		return
	}

	var image models.Image
	if err := database.DB.Where("id = ?", imageID).First(&image).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Erro ao buscar imagem: %s", err.Error()),
		})
		return
	}

	renderedHTML, err := renderImageBox(image)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": fmt.Sprintf("Erro ao renderizar imagem: %s", err.Error()),
		})
		return
	}
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(renderedHTML))
}

func renderImageBox(image models.Image) (string, error) {

	var tmpl, err = template.ParseFiles("templates/dish/box-image.html")
	if err != nil {
		fmt.Println("Error loading template:", err)
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, image); err != nil {
		fmt.Println("Error rendering template:", err)
		return "", err
	}
	return buf.String(), nil
}
