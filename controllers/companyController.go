package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CompanyPage(c *gin.Context) {
	fmt.Println("Rendering company page")
	c.HTML(http.StatusOK, "company.html", gin.H{})
}
