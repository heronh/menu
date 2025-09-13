package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func CompanyPage(c *gin.Context) {
	c.HTML(http.StatusOK, "company.html", gin.H{})
}
