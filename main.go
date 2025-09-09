package main

import (
	"log"
	"main/database"
	"main/initializers"
	"os"
	"path/filepath"

	"main/controllers"

	"github.com/gin-gonic/gin"
)

func init() {
	// Load .env file (if any)
	initializers.LoadEnvVariables()

	// Connect to DB
	database.Connect()

	// Run Migrations
	initializers.SyncDatabase()

	// Create Privileges
	initializers.CreatePrivileges()

	// Create Super User
	initializers.CreateSuperUser()
}

func main() {
	log.Println("Application setup complete. Starting application...")

	r := gin.Default()

	// Serve template files and it's subfolders from the 'templates' directory
	// load templates recursively
	files, err := loadTemplates("templates")
	if err != nil {
		log.Println(err)
	}
	r.LoadHTMLFiles(files...)

	// Serve static files (CSS) from the 'static' directory
	r.Static("/static", "./static")
	// Serve Bootstrap icons from the 'node_modules/bootstrap-icons' directory
	r.Static("/icons", "./static/bootstrap-icons")

	r.GET("/", controllers.WelcomePage)
	r.GET("/todo", controllers.TodoPage)

	r.GET("/register.html", controllers.RegisterPage)

	// read port in .env file and starts the server
	host_port := os.Getenv("HOST_PORT")
	if host_port == "" {
		host_port = "8080" // default port if not specified
	}
	r.Run(":" + host_port)

}

func loadTemplates(root string) (files []string, err error) {
	err = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fileInfo, err := os.Stat(path)
		if err != nil {
			return err
		}
		if fileInfo.IsDir() {
			if path != root {
				loadTemplates(path)
			}
		} else {
			files = append(files, path)
		}
		return err
	})
	return files, err
}
