package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WelcomePage(c *gin.Context) {
	c.HTML(http.StatusOK, "welcome.html", gin.H{
		"title": "Welcome to Menu Management System",
	})
}
