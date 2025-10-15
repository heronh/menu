package controllers

import (
	"bytes"
	"fmt"
	"main/database"
	"main/models"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadDishImage(c *gin.Context) {

	fmt.Println("Uploading single image...")
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
	imageID := c.Param("image_id")

	// Fetch the dish image to ensure it exists
	var dishImage models.Image
	if err := database.DB.Where("id = ?", imageID).First(&dishImage).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching dish image: %v", err)
		return
	}

	// Check if image is not used in any dish
	var count int64
	database.DB.Model(&models.Dish{}).Where("image_id = ?", imageID).Count(&count)
	if count > 0 {
		c.String(http.StatusBadRequest, "Cannot delete image: it is currently used in one or more dishes.")
		return
	}

	// Delete the dish image record
	if err := database.DB.Delete(&dishImage).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error deleting dish image: %v", err)
		return
	}

	// Delete the image file from the filesystem
	if err := os.Remove(dishImage.OriginalFileName); err != nil {
		fmt.Printf("Warning: could not delete image file: %v\n", err)
		// Not returning error to user since the DB record is already deleted
	}
	c.JSON(http.StatusOK, gin.H{"success": "Image deleted successfully"})
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

	renderedHTML := ""
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
		UniqueName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
		fmt.Println("Saving file to:", uploadPath+UniqueName)
		if err := c.SaveUploadedFile(file, "."+uploadPath+UniqueName); err != nil {
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
			UniqueName:       uploadPath + UniqueName,
			Storage:          "local",
			InsertedByID:     userID.(uint),
			CompanyID:        companyID.(uint),
			CreatedAt:        time.Now(),
		}
		fmt.Println("With path:", uploadPath+UniqueName)
		// Save record to database
		if err := database.DB.Create(&dishImage).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": fmt.Sprintf("Erro ao salvar o registro da imagem: %s", err.Error()),
			})
			return
		}

		// Render image box HTML
		renderedBox, err := renderImageBox(dishImage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": fmt.Sprintf("Erro ao renderizar a imagem: %s", err.Error()),
			})
			return
		}
		renderedHTML += renderedBox
		fmt.Println("Image uploaded and record created successfully:", file.Filename)
	}
	c.JSON(http.StatusOK, gin.H{"success": "Section created successfully", "renderedHTML": renderedHTML})
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
