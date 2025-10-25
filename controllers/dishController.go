package controllers

import (
	"fmt"
	"main/database"
	"main/models"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"

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
		"title":           "Adicionar novo prato",
		"Sections":        Sections,
		"Images":          Images,
		"CompanyID":       CompanyID,
		"UserID":          UserID,
		"TimeList":        TimeList,
		"DishName":        "Batata frita " + strconv.Itoa(rand.Intn(1000)),
		"DishDescription": "Deliciosas batatas fritas crocantes, perfeitas para acompanhar qualquer refeição.",
	})
}

func CreateDish(c *gin.Context) {

	fmt.Println("CreateDish called")

	companyIDStr := c.PostForm("CompanyID")
	var companyID uint
	if companyIDStr != "" {
		if cid, err := strconv.ParseUint(companyIDStr, 10, 64); err == nil {
			companyID = uint(cid)
		}
	}

	userIDStr := c.PostForm("UserID")
	var userID uint
	if userIDStr != "" {
		if uid, err := strconv.ParseUint(userIDStr, 10, 64); err == nil {
			userID = uint(uid)
		}
	}

	sectionIDStr := c.PostForm("SectionID")
	var sectionID uint
	if sectionIDStr != "" {
		if sid, err := strconv.ParseUint(sectionIDStr, 10, 64); err == nil {
			sectionID = uint(sid)
		}
	}

	ActiveStr := c.PostForm("Active")
	if ActiveStr == "" {
		ActiveStr = c.PostForm("active")
	}
	fmt.Printf("Received Active form value: '%s'\n", ActiveStr)
	var ActiveBool bool
	if ActiveStr == "on" || ActiveStr == "true" || ActiveStr == "1" || ActiveStr == "checked" {
		ActiveBool = true
	} else {
		ActiveBool = false
	}
	// create pointer for Active to match models.Dish.Active (*bool)
	activePtr := ActiveBool

	ShowPriceStr := c.PostForm("ShowPrice")
	var ShowPriceBool bool
	if ShowPriceStr == "on" || ShowPriceStr == "true" || ShowPriceStr == "1" || ShowPriceStr == "checked" {
		ShowPriceBool = true
	} else {
		ShowPriceBool = false
	}
	// create pointer for ShowPrice to match models.Dish.ShowPrice (*bool)
	showPricePtr := ShowPriceBool

	ShowDescriptionStr := c.PostForm("ShowDescription")
	var ShowDescriptionBool bool
	if ShowDescriptionStr == "on" || ShowDescriptionStr == "true" || ShowDescriptionStr == "1" || ShowDescriptionStr == "checked" {
		ShowDescriptionBool = true
	} else {
		ShowDescriptionBool = false
	}
	// create pointer for ShowDescription to match models.Dish.ShowDescription (*bool)
	showDescriptionPtr := ShowDescriptionBool
	// Set Availability from comma-separated string if provided
	availabilityStr := c.PostForm("Availability")
	var Availability []string
	if availabilityStr != "" {
		for _, v := range strings.Split(availabilityStr, ",") {
			v = strings.TrimSpace(v)
			if v != "" {
				Availability = append(Availability, v)
			}
		}
	}

	// Set WeekDays from comma-separated string if provided
	weekDaysStr := c.PostForm("WeekDays")
	var WeekDays []string
	if weekDaysStr != "" {
		for _, v := range strings.Split(weekDaysStr, ",") {
			v = strings.TrimSpace(v)
			if v != "" {
				WeekDays = append(WeekDays, v)
			}
		}
	}

	PriceStr := c.PostForm("Price")
	Price, err := strconv.ParseFloat(PriceStr, 64)
	if err != nil {
		c.String(http.StatusBadRequest, "Preço inválido: %v", err)
		return
	}

	dish := models.Dish{
		Name:            c.PostForm("Name"),
		CompanyID:       companyID,
		UserID:          userID,
		SectionID:       sectionID,
		Active:          &activePtr,
		Description:     c.PostForm("Description"),
		Price:           Price,
		ShowPrice:       &showPricePtr,
		Availability:    Availability,
		ShowDescription: &showDescriptionPtr,
		WeekDays:        WeekDays,
	}

	fmt.Printf("Dish fields:\n")
	fmt.Printf("ID: %v\n", dish.ID)
	fmt.Printf("Name: %v\n", dish.Name)
	fmt.Printf("Description: %v\n", dish.Description)
	fmt.Printf("Active (from form): %v\n", ActiveBool)
	fmt.Printf("Price: %v\n", dish.Price)
	fmt.Printf("SectionID: %v\n", dish.SectionID)
	fmt.Printf("CompanyID: %v\n", dish.CompanyID)
	fmt.Printf("CreatedAt: %v\n", dish.CreatedAt)
	fmt.Printf("UpdatedAt: %v\n", dish.UpdatedAt)
	fmt.Printf("Availability: %v\n", dish.Availability)
	fmt.Printf("WeekDays: %v\n", dish.WeekDays)

	// Force GORM to include the Active field in the INSERT so a false value
	// is not accidentally omitted and the DB default (true) applied.
	// Use Debug() to print the generated SQL so we can verify Active is included
	if err := database.DB.Debug().Select("Name", "Description", "Price", "SectionID", "CompanyID", "UserID", "Active", "Availability", "WeekDays", "ShowPrice", "ShowDescription").Create(&dish).Error; err != nil {
		c.String(http.StatusInternalServerError, "Error creating dish: %v", err)
		return
	}

	// Read back the created record to confirm what was stored in DB
	var created models.Dish
	if err := database.DB.Where("id = ?", dish.ID).First(&created).Error; err == nil {
		if created.Active != nil {
			fmt.Printf("Created dish Active in DB: %v\n", *created.Active)
		} else {
			fmt.Printf("Created dish Active in DB: <nil>\n")
		}
	} else {
		fmt.Printf("Could not read back created dish: %v\n", err)
	}
	fmt.Println("Dish created with ID:", dish.ID)

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
