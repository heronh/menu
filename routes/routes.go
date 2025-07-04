package routes

import (
	"net/http"

	"example.com/m/v2/auth"
	"example.com/m/v2/handlers"
	"example.com/m/v2/models"
	"github.com/gin-gonic/gin"
)

// SetupRouter configures the Gin router and defines the routes.
func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Serve static files (HTML, CSS, JS) from a "static" directory
	// This will be needed for index.html, company.html, and Bootstrap
	// The actual directory name can be "static", "public", or "templates/static"
	router.Static("/static", "./static") // Example: /static/css/bootstrap.min.css

	// Serve HTML templates from a "templates" directory
	router.LoadHTMLGlob("templates/*") // Example: index.html, company.html

	// Index page route - serves index.html
	router.GET("/", func(c *gin.Context) {
		// For now, just a JSON message. Later, this will serve index.html.
		// c.HTML(http.StatusOK, "index.html", nil)
		// For now, keeping the original JSON response for basic API check
		c.JSON(http.StatusOK, gin.H{"message": "Welcome! Please use /web/index.html to access the frontend."})
	})

	// Route to serve the actual index.html (temporary, until we decide on full SPA or server-rendered)
	router.GET("/web/index.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Welcome",
		})
	})
	// Route to serve company.html (temporary, needs auth)
	// This route itself doesn't need JWT middleware if the page is static,
	// but the DATA it loads via API calls WILL need auth.
	router.GET("/web/company.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "company.html", gin.H{
			"title": "Company Administration",
		})
	})

	// API v1 routes
	apiV1 := router.Group("/api/v1")
	{
		// Authentication routes
		authRoutes := apiV1.Group("/auth")
		{
			authRoutes.POST("/register", handlers.RegisterUserCompanyHandler)
			authRoutes.POST("/login", handlers.LoginHandler)
		}

		// Company routes
		// These routes require authentication
		companyRoutes := apiV1.Group("/companies")
		companyRoutes.Use(auth.JWTMiddleware()) // Apply JWT middleware to all /companies routes
		{
			// Route for the logged-in user's own company data
			// GET /api/v1/companies/my - uses companyId from JWT
			// This is an alternative to /api/v1/companies/:companyId for non-SU users
			companyRoutes.GET("/my", handlers.GetCompanyDataHandler)

			// Routes for specific company ID - primarily for SU or for clarity
			// GET /api/v1/companies/:companyId
			companyRoutes.GET("/:companyId", handlers.GetCompanyDataHandler)

			// PUT /api/v1/companies/:companyId
			// Authorization for who can update (SU or Manager of that company) is handled within the handler.
			companyRoutes.PUT("/:companyId", handlers.UpdateCompanyDataHandler)

			// Users associated with a company
			// GET /api/v1/companies/:companyId/users
			companyRoutes.GET("/:companyId/users", handlers.ListCompanyUsersHandler)

			// Categories, Dishes, Images for a company
			// GET /api/v1/companies/:companyId/categories
			companyRoutes.GET("/:companyId/categories", auth.Authorize(models.PrivilegeSuperAdministrator, models.PrivilegeManager, models.PrivilegeEmployee), handlers.ListCompanyCategoriesHandler)

			// GET /api/v1/companies/:companyId/dishes
			companyRoutes.GET("/:companyId/dishes", auth.Authorize(models.PrivilegeSuperAdministrator, models.PrivilegeManager, models.PrivilegeEmployee), handlers.ListCompanyDishesHandler)

			// GET /api/v1/companies/:companyId/images
			companyRoutes.GET("/:companyId/images", auth.Authorize(models.PrivilegeSuperAdministrator, models.PrivilegeManager, models.PrivilegeEmployee), handlers.ListCompanyImagesHandler)

			// Potentially routes for creating categories, dishes, images within a company context
			// e.g., POST /api/v1/companies/:companyId/dishes (requires Manager or SU)
		}

		// TODO: Add routes for Categories, Dishes, Images (CRUD operations)
		// These might be top-level like /api/v1/dishes or nested if always company-specific

		// TODO: Add routes for Messages, Logs, Todos
	}

	return router
}
