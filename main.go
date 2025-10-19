package main

import (
	"log"
	"main/database"
	"main/initializers"
	"os"
	"path/filepath"

	"main/controllers"
	"main/middleware"

	"github.com/gin-gonic/gin"
)

func init() {

	// Load .env file (if any)
	initializers.LoadEnvVariables()

	// Connect to DB
	database.Connect()

	// Reinicia o banco de dados
	// fmt.Println("Resetting database...")
	// initializers.ResetDatabase(database.DB, []any{
	// 	&models.Section{},
	// 	&models.Image{},
	// 	&models.Dish{},
	// 	&models.Todo{},
	// 	&models.Privilege{},
	// 	&models.Message{},
	// 	&models.User{},
	// 	&models.Company{},
	// })

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
	// Serve uploaded files (images) from the 'uploads' directory
	r.Static("/uploads", "./uploads")

	r.GET("/", controllers.WelcomePage)

	// Funções relativas as tarefas
	r.GET("/todo", controllers.TodoPage)
	r.POST("/todo", controllers.SaveTodo)
	r.POST("/todo_delete", controllers.DeleteTodo)
	r.POST("/todo_check", controllers.CheckTodo)
	r.POST("/todo_uncheck", controllers.UncheckTodo)

	// Funções relativas aos usuários
	r.POST("/is-email-available", controllers.IsEmailAvailable)
	r.GET("/register", controllers.RegisterPage)
	r.POST("/register", controllers.RegisterCompanyUser)
	r.GET("/login", controllers.LoginPage)
	r.POST("/login", controllers.LoginUser)
	r.GET("/logout", controllers.LogoutUser)

	// Funções relativas às empresas
	r.GET("/company", middleware.JWTAuthMiddleware(), controllers.CompanyPage)

	// Funções relativas aos pratos
	r.GET("/dishes/new", middleware.JWTAuthMiddleware(), controllers.NewDishPage)
	r.POST("/dishes/dish/new", middleware.JWTAuthMiddleware(), controllers.CreateDish)
	r.GET("/dishes/edit/:id", middleware.JWTAuthMiddleware(), controllers.EditDishPage)
	r.POST("/dishes/edit/:id", middleware.JWTAuthMiddleware(), controllers.UpdateDish)
	r.POST("/dishes/delete/:id", middleware.JWTAuthMiddleware(), controllers.DeleteDish)
	r.POST("/dishes/dish/validate", middleware.JWTAuthMiddleware(), controllers.ValidateDish)

	// Funções relativas as imagens dos pratos
	r.POST("/dishes/images/upload/:id", middleware.JWTAuthMiddleware(), controllers.UploadDishImage)
	r.DELETE("/dishes/images/delete/:image_id", middleware.JWTAuthMiddleware(), controllers.DeleteDishImage)
	r.POST("/dishes/images/upload-multiple-images", middleware.JWTAuthMiddleware(), controllers.UploadMultipleDishImages)
	r.POST("/dishes/images/list/:id", middleware.JWTAuthMiddleware(), controllers.ListDishImages)

	// Funções relativas as seções dos pratos
	r.POST("/sections/new", middleware.JWTAuthMiddleware(), controllers.CreateSection)
	r.DELETE("/sections/delete/:id", middleware.JWTAuthMiddleware(), controllers.DeleteSection)

	// Exemplo de componentes do tailwind
	r.StaticFile("/components", "templates/components.html")

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
