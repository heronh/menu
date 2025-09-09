package controllers

import (
	"fmt"
	"net/http"
	"time"

	"main/database"
	"main/models"

	"github.com/gin-gonic/gin"
)

func UncheckTodo(c *gin.Context) {
	if err := check_uncheck(false, c); err != nil {
		return
	}
}

func CheckTodo(c *gin.Context) {
	if err := check_uncheck(true, c); err != nil {
		return
	}
}

func check_uncheck(status bool, c *gin.Context) error {

	Id, err := parseId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return err
	}
	fmt.Println("Checking todo with id:", Id)

	var todo models.Todo
	if err := database.DB.Where("id = ?", Id).First(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find todo"})
		return err
	}
	todo.Completed = status
	todo.UpdatedAt = time.Now()

	if err := database.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update todo"})
		return err
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully changed status of todo"})
	return nil
}

func DeleteTodo(c *gin.Context) {

	Id, err := parseId(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Deleting todo with id:", Id)
	if err := database.DB.Delete(&models.Todo{}, Id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not delete todo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Successfully deleted todo"})
}

func parseId(c *gin.Context) (int, error) {
	type RequestData struct {
		Id int `json:"Id"`
	}
	var requestData RequestData
	if err := c.BindJSON(&requestData); err != nil {
		return 0, err
	}
	return requestData.Id, nil
}

func SaveTodo(c *gin.Context) {
	fmt.Println("Creating todo")
	var todo models.Todo
	todo.CreatedAt = time.Now()
	todo.UpdatedAt = time.Now()
	todo.Completed = false
	todo.Description = c.PostForm("description")

	fmt.Println(c)

	// look for first user by privilege
	var user models.User
	user, err := getUserByPrivilege("su")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find user with privilege 'su'"})
		return
	}
	fmt.Println("User found:", user.Email)

	todo.UserID = user.ID
	if err := database.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create todo"})
		return
	}
	c.Redirect(http.StatusFound, "/todo")
}

func TodoPage(c *gin.Context) {

	fmt.Println("Retrieving todos")
	var todos []models.Todo
	if err := database.DB.Order("completed, updated_at desc").Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve todos"})
		return
	}
	for _, todo := range todos {
		fmt.Printf("Todo: %s\n", todo.Description)
	}

	// look for first user by privilege
	var user models.User
	user, err := getUserByPrivilege("su")
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Could not find user with privilege 'su'"})
		return
	}
	fmt.Println("User found:", user.Email)

	// retrieve email and user id from the context
	Email := user.Email
	ID := user.ID
	c.HTML(http.StatusOK, "todo.html", gin.H{
		"Todos": todos,
		"Email": Email,
		"Id":    ID,
	})
}

func getUserByPrivilege(slug string) (models.User, error) {

	var user models.User

	// List all available privileges
	var privilege models.Privilege
	if err := database.DB.First(&privilege, "slug = ?", slug).Error; err != nil {
		return user, err
	}
	fmt.Println("Privilege found:", privilege.Name, "Slug:", privilege.Slug)
	// Find first user with this privilege
	if err := database.DB.First(&user, "privilege_id = ?", privilege.ID).Error; err != nil {
		return user, err
	}
	return user, nil
}
