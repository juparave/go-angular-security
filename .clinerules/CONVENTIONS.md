# Go Fiber + Gorm Project Conventions

This document outlines common conventions and best practices for developing web applications using Go with the Fiber framework and
Gorm ORM. Adhering to these conventions promotes consistency, readability, and maintainability across the project.

## 1. Project Structure

Follow the standard Go project layout:

```
├── cmd/
│   └── api/
│       └── main.go         # Main application entry point
├── internal/
│   ├── config/             # Configuration loading
│   ├── database/           # Database connection setup (Gorm)
│   ├── handlers/           # Fiber request handlers (controllers)
│   ├── middleware/         # Fiber middleware
│   ├── models/             # Gorm data models (structs)
│   ├── routes/             # Fiber route definitions
│   ├── services/           # Business logic services
│   └── utils/              # Utility functions
├── pkg/                    # Reusable libraries
├── migrations/             # Database migration files (if using external tool)
├── go.mod
├── go.sum
└── ref/
    └── CONVENTIONS.md      # This file
```

## 2. Naming Conventions

-   **Packages:** `lowercase`, short, and descriptive (e.g., `handlers`, `models`, `database`).
-   **Files:** `snake_case.go` (e.g., `user_handler.go`, `db_connect.go`).
-   **Variables:** `camelCase`. Exported variables start with an uppercase letter (`PascalCase`).
-   **Constants:** `PascalCase` for exported constants, `camelCase` or `UPPER_SNAKE_CASE` for unexported, depending on context and
team preference.
-   **Functions/Methods:** `camelCase` for unexported, `PascalCase` for exported.
-   **Structs:** `PascalCase`.
-   **Interfaces:** `PascalCase`, often ending with `-er` (e.g., `Reader`, `Writer`) if they represent an action.
-   **Database Tables:** `snake_case`, plural (e.g., `users`, `product_orders`).
-   **Database Columns:** `snake_case` (e.g., `first_name`, `created_at`).
-   **Gorm Models:** `PascalCase`, singular (e.g., `User`, `ProductOrder`). Match struct names to table names (Gorm handles
pluralization/snake_case conversion, but be explicit if needed).

## 3. Fiber Conventions

-   **Handlers:**
    -   Signature: `func HandlerName(c *fiber.Ctx) error`. Handlers should return an `error`. Fiber handles nil errors as success
(often `200 OK` implicitly or `204 No Content` for methods like `DELETE`).
    -   Return non-nil errors for client or server issues. Use Fiber's error handler or custom middleware to map errors to HTTP
responses.
    -   Keep handlers thin; delegate business logic to `services`.
-   **Responses:**
    -   Use net/http status constants for consistency
    -   Use `c.Status(http.StatusOK).JSON(data)` for successful JSON responses.
    -   Use `c.SendStatus(http.StatusNoContent)` for successful responses with no body.
    -   Use `c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "message"})` for client errors.
-   **Error Handling:**
    -   Return specific error types or use `fiber.NewError(statusCode, message)` for standard HTTP errors.
    -   Implement a custom global error handler (`fiber.Config{ErrorHandler: ...}`) to centralize error response formatting.
-   **Middleware:**
    -   Place common logic (auth, logging, CORS, rate limiting) in middleware.
    -   Organize middleware in the `internal/middleware` package.
-   **Routing:**
    -   Define routes in `internal/routes/routes.go` or similar.
    -   Group related routes using `app.Group("/path", optionalMiddleware)`.
-   **Input Binding/Validation:**
    -   Use `c.BodyParser(&dto)` to parse request bodies into structs.
    -   Use a validation library (like `go-playground/validator`) integrated via struct tags or middleware for input validation.
Return `400 Bad Request` on validation failure.

## 4. Gorm Conventions

-   **Models:**
    -   Define structs in `internal/models`.
    -   Use `gorm:"..."` tags for column names, types, constraints (e.g.,
`gorm:"column:user_name;type:varchar(100);uniqueIndex"`).
    -   Use `json:"..."` tags for JSON serialization control (e.g., `json:"firstName"`, `json:"-"` to omit).
    -   Embed `gorm.Model` for standard fields (`ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`) or define them manually if needed.
-   **Database Connection:**
    -   Initialize the Gorm DB instance in `internal/database` and make it accessible (e.g., via a global variable or dependency
injection).
    -   Handle connection errors gracefully during startup.
-   **CRUD Operations:**
    -   Use standard Gorm methods: `db.Create()`, `db.First()`, `db.Find()`, `db.Where()`, `db.Save()`, `db.Updates()`,
`db.Delete()`.
    -   Always check for errors returned by Gorm operations.
-   **Error Handling:**
    -   Specifically check for `gorm.ErrRecordNotFound` when expecting a single result (e.g., `db.First()`).
    -   Propagate other database errors up the call stack to be handled appropriately (e.g., logged and potentially returned as a
`500 Internal Server Error`).
-   **Transactions:**
    -   Use `db.Transaction(func(tx *gorm.DB) error { ... })` for operations requiring atomicity.
-   **Migrations:**
    -   Use Gorm's `AutoMigrate` for simple schema management during development.
    -   For production, prefer dedicated migration tools (e.g., `golang-migrate/migrate`, `Atlas`, `Gormigrate`) and store
migration files (often SQL) in the `migrations/` directory.
-   **Preloading (Eager Loading):**
    -   Use `db.Preload("AssociationName").Find(&results)` to avoid N+1 query problems.

## 5. General Go Conventions

-   **Error Handling:**
    -   Use `errors.New` or `fmt.Errorf` for simple errors.
    -   Wrap errors using `fmt.Errorf("context: %w", err)` to add context while preserving the original error type (use with
`%w`).
    -   Check errors immediately after a function call that returns one (`if err != nil { ... }`).
-   **Formatting:** Use `gofmt` or `goimports` to automatically format code before committing.
-   **Testing:** Write unit tests (`_test.go` files) for handlers, services, and utilities. Consider integration tests that
involve the database and potentially HTTP requests.
-   **Configuration:** Load configuration from environment variables (using libraries like `viper` or `godotenv`) or configuration
files. Avoid hardcoding configuration values.
-   **Logging:** Use a structured logging library (e.g., `slog` (Go 1.21+), `zap`, `zerolog`) for application logging.
