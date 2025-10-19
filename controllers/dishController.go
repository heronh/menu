package controllers

import (
	"errors"
	"fmt"
	"main/database"
	"main/models"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Estrutura para formatar a resposta de erro
type FieldError struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Msg   string `json:"message"` // Mensagem de erro customizada (opcional)
}

func getErrorMsg(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "Este campo é obrigatório"
	case "min":
		return "O valor é muito curto"
	case "gt":
		return "Deve ser maior que zero"
	// Adicione mais casos conforme necessário
	default:
		return "Erro de validação no campo"
	}
}

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
	CompanyID, _ := c.Get("company_id")
	UserID, _ := c.Get("user_id")
	fmt.Println("CompanyID: ", CompanyID)
	fmt.Println("UserID: ", UserID)

	Sections := []models.Section{}
	if err := database.DB.Where("company_id = ?", CompanyID).Find(&Sections).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching sections: %v", err)
		return
	}
	fmt.Println("Total sections found:", len(Sections))

	var Images []models.Image
	if err := database.DB.Where("company_id = ?", CompanyID).Find(&Images).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error fetching images: %v", err)
		return
	}
	fmt.Println("Total images found:", len(Images))

	c.HTML(http.StatusOK, "new-dish.html", gin.H{
		"title":     "Adicionar novo prato",
		"Sections":  Sections,
		"Images":    Images,
		"CompanyID": CompanyID,
		"UserID":    UserID,
		"TimeList":  TimeList,
	})
}

func CreateDish(c *gin.Context) {

	fmt.Println("CreateDish called")
	dish := models.Dish{}
	if err := c.ShouldBind(&dish); err != nil {

		// 1. Tenta fazer type assertion para ValidationErrors (erros de validação)
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			// Há erros de validação, podemos iterar sobre eles
			out := make([]FieldError, len(ve))
			for i, fe := range ve {
				out[i] = FieldError{
					Field: fe.Field(), // Nome do campo na struct (ex: Name, Price)
					Tag:   fe.Tag(),   // Tag de validação que falhou (ex: required, min)
					Msg:   getErrorMsg(fe),
				}
			}

			// Retorna uma resposta JSON detalhada com os erros
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Erros de validação nos campos",
				"errors":  out,
			})
			c.String(http.StatusInternalServerError, "Error parsing dish: %v", err)
			return
		}

		// 2. Se não for um erro de validação (ex: JSON malformado, tipo de dado incorreto)
		// Retorna o erro genérico de binding
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Erro no formato dos dados ou binding",
			"details": err.Error(),
		})
		c.String(http.StatusInternalServerError, "Error dish: %v", err)
		return
	}

	dish.CreatedAt = time.Now()
	dish.UpdatedAt = time.Now()
	fmt.Printf("Dish fields:\n")
	fmt.Printf("ID: %v\n", dish.ID)
	fmt.Printf("Name: %v\n", dish.Name)
	fmt.Printf("Description: %v\n", dish.Description)
	fmt.Printf("Price: %v\n", dish.Price)
	fmt.Printf("SectionID: %v\n", dish.SectionID)
	fmt.Printf("CompanyID: %v\n", dish.CompanyID)
	fmt.Printf("CreatedAt: %v\n", dish.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", dish.UpdatedAt)

	newDish := models.Dish{
		Name:        dish.Name,
		Description: dish.Description,
		Price:       dish.Price,
		SectionID:   dish.SectionID,
		CompanyID:   dish.CompanyID,
		UserID:      dish.UserID,
	}
	if err := database.DB.Create(&newDish).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error creating dish: %v", err)
		return
	}
	fmt.Println("Dish created with ID:", newDish.ID)

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
	dish := models.Dish{
		Name: c.PostForm("name"),
		SectionID: func() uint {
			sectionIDStr := c.PostForm("section_id")
			sectionID, err := strconv.ParseUint(sectionIDStr, 10, 64)
			if err != nil {
				return 0
			}
			return uint(sectionID)
		}(),
	}
	fmt.Println("Validating dish:", dish.Name, "SectionID:", dish.SectionID)

	name := ""
	if dish.Name == "" {
		name = "Nome do prato é obrigatório."
	} else if len(dish.Name) < 5 {
		name = "Nome do prato deve ter pelo menos 5 caracteres."
	}

	section_id := ""
	if dish.SectionID == 0 {
		section_id = "Seção do prato é obrigatória."
	}

	valid := true
	if name != "" || section_id != "" {
		valid = false
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": valid,
		"errors": gin.H{
			"name":       name,
			"section_id": section_id,
		},
	})
}
