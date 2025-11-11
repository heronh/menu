package controllers

import (
	"fmt"
	"main/database"
	"main/models"
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
)

func getTranslateWeekdayToPortuguese(weekday string) string {

	wd := time.Now().Weekday()
	var currentWeekday string
	switch wd {
	case time.Monday:
		currentWeekday = "Segunda"
	case time.Tuesday:
		currentWeekday = "Terça"
	case time.Wednesday:
		currentWeekday = "Quarta"
	case time.Thursday:
		currentWeekday = "Quinta"
	case time.Friday:
		currentWeekday = "Sexta"
	case time.Saturday:
		currentWeekday = "Sábado"
	case time.Sunday:
		currentWeekday = "Domingo"
	default:
		currentWeekday = wd.String()
	}
	return currentWeekday
}

func IsDishAvailableNow(dish *models.Dish) bool {

	// Check availability based on dish.Availability and dish.WeekDays
	Availability := false
	// treat empty availability slice as available and consider "all" as wildcard
	if len(dish.Availability) == 0 {
		Availability = true
	} else {
		if slices.Contains(dish.Availability, "all") {
			Availability = true
		}
	}

	if !Availability {
		// Catch current time in "HH:MM" format
		currentTime := time.Now().Format("15:04")
		for _, timeRange := range dish.Availability {
			var startTime, endTime string
			fmt.Sscanf(timeRange, "%5s - %5s", &startTime, &endTime)
			fmt.Println("currentTime:", currentTime, ", startTime:", startTime, ", endTime:", endTime)
			if currentTime >= startTime && currentTime <= endTime {
				Availability = true
				break
			}
		}
	}

	WeekDays := false
	// treat empty weekdays slice as available and consider "all" as wildcard
	if len(dish.WeekDays) == 0 {
		WeekDays = true
	} else {
		if slices.Contains(dish.WeekDays, "all") {
			WeekDays = true
		}
	}

	if !WeekDays {
		currentWeekday := getTranslateWeekdayToPortuguese(time.Now().Weekday().String())
		if slices.Contains(dish.WeekDays, currentWeekday) {
			WeekDays = true
		}
	}
	fmt.Println("Availability:", Availability, "WeekDays:", WeekDays)

	return Availability && WeekDays
}

func ViewCompanyMenu(c *gin.Context) {

	// Get all dishes from this company
	CompanyID, _ := c.Get("company_id")
	var dishes []models.Dish
	if err := database.DB.
		Preload("Images").
		Preload("Section").
		Where("company_id = ?", CompanyID).
		Order("created_at DESC").
		Find(&dishes).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dishes"})
		return
	}

	for i := range dishes {
		fmt.Println("\n", dishes[i].Name)
		//fmt.Println("Active:", dishes[i].Active)
		//fmt.Println(dishes[i].Description)
		//fmt.Println(dishes[i].Price)
		//fmt.Println("Section: ", dishes[0].Section.Description)
		fmt.Println("Availability:", dishes[i].Availability)
		fmt.Println("WeekDays:", dishes[i].WeekDays)
		dishes[i].DishSection = dishes[i].Section.Description
		fmt.Println("Dish Section:", dishes[i].DishSection)
		dishes[i].AvailableNow = IsDishAvailableNow(&dishes[i])
		fmt.Println("Available Now:", dishes[i].AvailableNow)

		fmt.Println("Images:")
		// You can access dish.Images here
		_ = dishes[i].Images // Just to avoid unused variable warning
		for _, img := range dishes[i].Images {
			fmt.Println(" - Image URL:", img.OriginalFileName)
			fmt.Println(" - Stored as:", img.UniqueName)
		}
	}

	// Render the menu template
	c.HTML(http.StatusOK, "menu-company.html", gin.H{
		"Title":     "Menu",
		"Dishes":    dishes,
		"CompanyID": CompanyID,
	})
}

func ViewCustomerMenu(c *gin.Context) {

	CompanyID, _ := c.Get("company_id")

	// Get all dishes from this company
	var dishes []models.Dish
	if err := database.DB.Preload("Images").Where("company_id = ?", CompanyID).Find(&dishes).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dishes"})
		return
	}

	// Render the menu template
	c.HTML(http.StatusOK, "menu.html", gin.H{
		"Title":     "Menu",
		"Dishes":    dishes,
		"CompanyID": CompanyID,
	})
}
