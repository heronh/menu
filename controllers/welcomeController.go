package controllers

import (
	"github.com/gin-gonic/gin"
)

func WelcomePage(c *gin.Context) {
	c.HTML(200, "welcome.html", gin.H{
		"title": "Welcome to Menu Management System",
	})
}

