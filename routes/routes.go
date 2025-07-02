package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures the Gin router and defines the routes.
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Index page route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Gin GORM Application!"})
	})

	// You can add more routes here for your models
	// e.g., /users, /companies, etc.

	return router
}
