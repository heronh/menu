package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterPage(c *gin.Context) {
	fmt.Println("Acessou a página de registro")
	c.HTML(http.StatusOK, "register.html",
		gin.H{
			"Title": "Cadastro de Usuário",
		})
}
