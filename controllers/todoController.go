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
	Id := c.PostForm("Id")
	var userModel models.User
	if err := database.DB.Where("id = ?", Id).First(&userModel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not find user"})
		return
	}

	if err := database.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create todo"})
		return
	}
	c.Redirect(http.StatusFound, "/todo")
}

func TodoPage(c *gin.Context) {

	var todos []models.Todo
	if err := database.DB.Order("completed, updated_at desc").Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve todos"})
		return
	}

	// retrieve email and user id from the context
	Email, _ := c.Get("email")
	ID, _ := c.Get("ID")
	c.HTML(http.StatusOK, "todo.html", gin.H{
		"Todos": todos,
		"Email": Email,
		"Id":    ID,
	})
}
