# Go Backend API for Secure Web Application

This directory (`server/`) contains the Go backend API for the secure web application. It's built using the [Fiber](https://docs.gofiber.io/) web framework and provides RESTful endpoints for the Angular frontend.

## Project Overview

The Go backend is responsible for:
- User authentication (registration, login) and JWT generation/validation.
- Authorization and role-based access control (RBAC).
- Managing user data, roles, and permissions.
- Serving data to the frontend and handling business logic.
- Securely interacting with a database (e.g., PostgreSQL, MySQL - specific configuration in `database/connect.go`).

## Prerequisites

- Go (version 1.19 or higher)
- Access to a database supported by GORM (the project is set up for PostgreSQL, but can be adapted).

## Setup and Running

1.  **Navigate to the server directory:**
    ```bash
    cd server
    ```

2.  **Install Dependencies:**
    If this is the first time or if `go.mod` or `go.sum` has changed:
    ```bash
    go mod tidy
    go mod download
    ```

3.  **Database Configuration:**
    - The database connection is configured in `database/connect.go`.
    - Ensure you have a running database instance (e.g., PostgreSQL).
    - Update the DSN (Data Source Name) string in `database/connect.go` with your database credentials:
      ```go
      dsn := "host=localhost user=youruser password=yourpassword dbname=yourdbname port=5432 sslmode=disable TimeZone=Asia/Shanghai"
      // Adjust according to your database setup
      ```
    - The application will attempt to automatically migrate the database schema upon startup (defined in `database/connect.go` using `AutoMigrate`).

4.  **Environment Variables (Optional but Recommended):**
    For sensitive information like JWT secrets or database credentials, it's recommended to use environment variables instead of hardcoding. This project might require you to set up:
    - `JWT_SECRET`: A secret key for signing JWTs.
    - `DB_HOST`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_PORT`: For database connection.
    *(Note: Current implementation might have these hardcoded; refactoring to use environment variables is a good practice.)*

5.  **Run the Application:**
    - **Without auto-reload:**
      ```bash
      go run main.go
      ```
    - **With auto-reload using Air (recommended for development):**
      Ensure [Air](https://github.com/cosmtrek/air) is installed (`go install github.com/cosmtrek/air@latest`).
      ```bash
      air
      # Air uses the .air.toml configuration file in this directory
      ```
    The API server will start, typically on `http://localhost:3000`.

## API Endpoints

The API routes are defined in `routes/routes.go`. Key endpoint groups include:

*   **Auth Endpoints (`/api/auth`):**
    *   `POST /api/auth/register`: User registration.
    *   `POST /api/auth/login`: User login, returns a JWT.
    *   `GET /api/auth/user`: Get current authenticated user's details (requires JWT).
    *   `POST /api/auth/logout`: User logout (typically invalidates session/cookie).
    *   `PUT /api/auth/users/info`: Update authenticated user's info.
    *   `PUT /api/auth/users/password`: Update authenticated user's password.

*   **User Management Endpoints (`/api/users` - often admin-restricted):**
    *   `GET /api/users`: List all users.
    *   `POST /api/users`: Create a new user.
    *   `GET /api/users/:id`: Get a specific user by ID.
    *   `PUT /api/users/:id`: Update a specific user by ID.
    *   `DELETE /api/users/:id`: Delete a specific user by ID.

*   **Role Management Endpoints (`/api/roles` - often admin-restricted):**
    *   `GET /api/roles`: List all roles.
    *   `POST /api/roles`: Create a new role.
    *   `GET /api/roles/:id`: Get a specific role by ID.
    *   `PUT /api/roles/:id`: Update a specific role by ID.
    *   `DELETE /api/roles/:id`: Delete a specific role by ID.

*   **Permission Management Endpoints (`/api/permissions` - often admin-restricted):**
    *   `GET /api/permissions`: List all permissions.

*(Refer to `routes/routes.go` and controller functions in `controllers/` for detailed request/response formats and specific logic.)*

## Database Schema

The database schema is defined using GORM models in the `models/` directory. Key models include:

-   **`User` (`models/users.go`):** Stores user information, credentials, and associated role.
    - Fields: `ID`, `FirstName`, `LastName`, `Email`, `Password`, `Phone`, `RoleID`, `Role`.
-   **`Role` (`models/role.go`):** Defines user roles (e.g., admin, editor, viewer).
    - Fields: `ID`, `Name`, `Permissions`.
-   **`Permission` (`models/permission.go`):** Defines granular permissions that can be assigned to roles.
    - Fields: `ID`, `Name`.
-   **`PasswordReset` (`models/users.go`):** (If password reset functionality is implemented)
    - Typically stores a token and expiry for password reset requests.

The database connection logic in `database/connect.go` includes an `AutoMigrate` call, which attempts to create or update the tables based on these models when the application starts.

## Authentication and Authorization

*   **Authentication:**
    - Implemented using JSON Web Tokens (JWT).
    - Upon successful login (`/api/auth/login`), a JWT is generated and sent to the client (typically as an HTTP-only cookie or in the response body).
    - The JWT secret key is crucial for security and should be managed carefully (ideally via environment variables). Current JWT logic is in `util/jwt.go`.
    - The `AuthMiddleware` (`middleware/authMiddleware.go`) protects routes that require authentication by validating the JWT from incoming requests.

*   **Authorization (Role-Based Access Control - RBAC):**
    - Users are assigned roles, and roles are associated with specific permissions.
    - The `PermissionMiddleware` (`middleware/permissionMiddleware.go`) can be used to protect endpoints based on the permissions required.
    - It checks if the authenticated user's role has the necessary permission to access a resource or perform an action.
    - Permissions are typically strings like "view_users", "edit_users", "view_products", etc. These are checked against the permissions assigned to the user's role.

## Project Structure

-   **`main.go`**: Entry point of the application. Initializes Fiber, database, and routes.
-   **`controllers/`**: Contains handler functions for API endpoints (business logic).
-   **`database/`**: Database connection and migration logic.
-   **`middleware/`**: Custom middleware (e.g., authentication, authorization).
-   **`models/`**: GORM models defining the database schema.
-   **`routes/`**: API route definitions.
-   **`util/`**: Utility functions (e.g., JWT handling, password hashing, input validation).
-   **`.air.toml`**: Configuration file for the Air live-reload tool.
-   **`go.mod`, `go.sum`**: Go module files for dependency management.

## Testing

Basic HTTP tests can be found in `controllers/test.http` which can be used with VSCode's REST Client extension.
Unit tests for utility functions like crypto can be found in `util/crypto_test.go`. To run tests:
```bash
go test ./...
```

Further contributions should include relevant tests for new functionalities.
---

This README provides a starting point. Refer to the source code for the most up-to-date details.
