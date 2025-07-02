# Restaurant Menu Web Application

## Description

This project is a Golang-based web server designed to display restaurant menus to consumers. It provides a seamless experience for browsing menu items, placing orders, and managing the restaurant's offerings from an administrative backend. PostgreSQL is used for robust data storage.

## Features

### Consumer-Facing Pages

*   **Homepage:** Displays the restaurant menu with categories and items.
*   **Item Details Page:** Shows detailed information about a specific menu item.
*   **Cart Page:** Displays items selected by the user and allows for order finalization.
*   **Checkout Page:** Enables users to enter payment information and complete their order.
*   **Confirmation Page:** Shows a confirmation of the placed order.

### Administrative Pages

*   **Login Page:** Allows administrators to log into the system.
*   **Menu Management Page:** Enables administrators to add, edit, or remove menu items.
*   **Order Management Page:** Allows administrators to view and manage customer orders.
*   **Reports Page:** Displays reports on sales, best-selling items, and other metrics.
*   **Settings Page:** Allows administrators to adjust system settings, such as operating hours and accepted payment methods.

## User Roles

*   **Super Administrator:** Has access to all administrative pages and can manage all aspects of the system, including users, permissions, and advanced settings.
*   **Administrator:** Has access to specific administrative pages, such as menu and order management, but will not have access to advanced settings or user management.

## Technology Stack

*   **Backend:** Golang, Gin Framework
*   **Database:** PostgreSQL
*   **Frontend:** HTML, CSS, JavaScript (with a responsive design for mobile-friendliness)

## Getting Started

(Placeholder for instructions on how to set up, configure, and run the project locally. This section will be updated as the project develops.)

## Contributing

(Placeholder for guidelines on how to contribute to the project. This section will be updated as the project develops.)

## License

(Placeholder for project license information. Consider using a standard open-source license like MIT or Apache 2.0.)
=======
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

