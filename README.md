# Gin GORM PostgreSQL Application

This is a sample Go application demonstrating the use of the Gin framework for building web APIs and GORM for database interaction with PostgreSQL.

## Features

- Basic project structure for a Go web application.
- Gin for HTTP routing and request handling.
- GORM for ORM (Object-Relational Mapping) with PostgreSQL.
- Data models: User, Company, Role, Dish, DishCategory, Image.
- Database auto-migration.
- A simple welcome message on the index page (`/`).

## Prerequisites

- Go (version 1.18 or higher recommended)
- PostgreSQL server running

## Setup

1.  **Clone the repository (if applicable) or download the files.**

2.  **Set up PostgreSQL:**
    *   Ensure your PostgreSQL server is running.
    *   Create a database (e.g., `gin_gorm_app`).
    *   You can use the default connection string in `database/database.go` or set the `DATABASE_URL` environment variable.
        Example `DATABASE_URL`:
        ```
        export DATABASE_URL="host=localhost user=youruser password=yourpassword dbname=gin_gorm_app port=5432 sslmode=disable TimeZone=UTC"
        ```
        The default DSN used if `DATABASE_URL` is not set is:
        `host=localhost user=postgres password=postgres dbname=gin_gorm_app port=5432 sslmode=disable TimeZone=Asia/Shanghai`
        **Note:** You might need to adjust the `dbname`, `user`, and `password` to match your PostgreSQL setup.

3.  **Install Go dependencies:**
    Navigate to the project root directory and run:
    ```bash
    go mod tidy
    # or
    go get .
    ```

## Running the Application

1.  **Run the application:**
    ```bash
    go run main.go
    ```
    The server will start, typically on port `8080`. You should see log messages indicating successful database connection, migration, and server startup.

2.  **Access the application:**
    Open your web browser or use a tool like `curl` to access the index page:
    ```bash
    curl http://localhost:8080/
    ```
    You should receive a JSON response:
    ```json
    {"message":"Welcome to the Gin GORM Application!"}
    ```

## Project Structure

```
.
├── go.mod
├── go.sum
├── main.go           // Main application entry point
├── database/
│   └── database.go   // Database connection and migration
├── models/
│   └── models.go     // GORM data models
├── routes/
│   └── routes.go     // Gin router and route definitions
└── README.md
```

## Further Development

-   Add more specific routes and handlers for each model (CRUD operations).
-   Implement authentication and authorization.
-   Write unit and integration tests.
-   Add request validation.
-   Structure handlers into their own package (e.g., `handlers/`).
