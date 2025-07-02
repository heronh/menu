## Agent Instructions for Gin GORM PostgreSQL Application

Welcome, agent! Here are some guidelines for working with this codebase.

### 1. Code Conventions

*   **Go Standard Practices:** Follow standard Go formatting (`gofmt`) and linting (`golint` or preferred linter).
*   **Error Handling:** Handle errors explicitly. Avoid panicking unless absolutely necessary (e.g., during initial setup if a critical component fails). Log errors appropriately.
*   **Modularity:** Keep packages focused.
    *   `models/`: Contains GORM data structures.
    *   `database/`: Handles database connection and GORM instance.
    *   `routes/`: Defines Gin HTTP routes.
    *   `handlers/` (if created): Contains Gin handler functions for business logic.
    *   `main.go`: Application entry point, initialization.
*   **Configuration:** Prefer environment variables for configuration (e.g., database DSN, server port). Provide sensible defaults if an environment variable is not set.
*   **Comments:** Write clear and concise comments for public functions, structs, and complex logic.

### 2. Database Migrations

*   GORM's `AutoMigrate` is used in `database/database.go`.
*   `AutoMigrate` will only create tables, add missing columns, and create missing indexes. It **will not** change existing column types or delete unused columns.
*   For more complex schema changes (e.g., renaming columns, changing types with data transformation), you will need to implement a proper migration system (e.g., using `golang-migrate/migrate` or GORM's migration tools if more advanced features are needed). For now, `AutoMigrate` is sufficient for initial development.

### 3. Dependencies

*   Dependencies are managed using Go Modules (`go.mod` and `go.sum`).
*   When adding or updating dependencies, run `go mod tidy`.

### 4. Running the Application

*   The application can be run using `go run main.go`.
*   Ensure PostgreSQL is running and accessible. The connection DSN can be configured via the `DATABASE_URL` environment variable or defaults to a common local setup (see `database/database.go` or `README.md`).

### 5. Testing

*   (No tests are currently implemented in this initial setup)
*   When adding tests:
    *   Place unit tests in the same package as the code they test, using the `_test.go` suffix (e.g., `user_test.go`).
    *   Consider using a separate database for testing or transaction-based rollbacks to keep tests isolated.
    *   For API endpoint testing, use the `net/http/httptest` package.

### 6. Adding New Features (e.g., CRUD for models)

*   **Define Model (if new):** Add/update struct in `models/models.go`.
*   **Update Migration:** Ensure `AutoMigrate` in `database/database.go` includes the new/updated model.
*   **Create Handlers:**
    *   Create a new file in a `handlers/` directory (e.g., `handlers/user_handlers.go`).
    *   Implement Gin handler functions for CRUD operations (Create, Read, Update, Delete). These handlers will contain the business logic (interacting with the database via GORM).
*   **Define Routes:** Add new routes in `routes/routes.go` and map them to your handler functions.
*   **Update `main.go` (if necessary):** Usually not needed unless new packages require initialization.

### 7. Logging

*   Use the standard `log` package for simple logging. For more advanced logging, consider a structured logging library (e.g., `logrus`, `zap`).

### 8. Environment Variables

*   `DATABASE_URL`: PostgreSQL connection string.
*   `PORT`: (Not implemented yet, but a good addition) Port for the Gin server to listen on.

Remember to keep the `README.md` updated with any significant changes to setup or operation. Good luck!
