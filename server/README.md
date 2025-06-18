# Go Backend Server

This Go backend provides a RESTful API for the Angular frontend. It handles business logic, data persistence, user authentication, and authorization.

## Key Technologies

*   **Go:** Version 1.24.2 (as per `go.mod`)
*   **Framework:** [Fiber v2](https://docs.gofiber.io/) - A fast and expressive web framework.
*   **ORM:** [GORM](https://gorm.io/) - A developer-friendly ORM for Go.
*   **Database:** [SQLite](https://www.sqlite.org/) - Used via GORM for data storage. The database file is typically created automatically.
*   **Authentication:** JWT (JSON Web Tokens) for secure API access.

## Project Structure

The `server/` directory is organized as follows:

*   `cmd/api/main.go`: Main entry point for the application.
*   `internal/`: Contains the core application logic, separated into:
    *   `config/`: Configuration management (currently minimal, likely relies on environment variables).
    *   `database/`: Database connection and GORM setup (`connect.go`).
    *   `handlers/`: HTTP request handlers (controllers) that process incoming requests.
    *   `middleware/`: Custom middleware for request processing (e.g., authentication, permissions).
    *   `models/`: GORM data models and entity definitions.
    *   `routes/`: API route definitions.
    *   `services/`: Business logic and service layer.
    *   `utils/`: Utility functions (e.g., for cryptography, JWT handling).
*   `go.mod`, `go.sum`: Go module files defining project dependencies.
*   `.air.toml`: Configuration file for `air` live reloading.
*   `api_spec.md`: Detailed API endpoint documentation.

## Setup and Configuration

1.  **Install Go:** Ensure you have Go (version 1.21 or newer recommended, project uses 1.24.2) installed on your system.
2.  **Navigate to Server Directory:**
    ```bash
    cd server
    ```
3.  **Install Dependencies:**
    Download and install the necessary Go modules:
    ```bash
    go mod download
    ```
    You can also use `go mod tidy` to ensure your `go.mod` file matches your source code.
4.  **Configuration:**
    *   The application is configured to run with default settings suitable for local development.
    *   Port: The server listens on port `5000` by default (see `cmd/api/main.go`).
    *   CORS: Configured to allow requests from `http://localhost:4200` (the default Angular development server).
    *   Database: Uses SQLite, and the database file (`gorm.db` or similar, as per GORM defaults) is typically created automatically in the `server` directory when the application starts. No manual database setup is usually required for development.
    *   Environment Variables: While not explicitly detailed in a config file yet, be aware that future configurations or production deployments might rely on environment variables.

## Running the Server

There are two main ways to run the server:

### 1. With Live Reloading (using `air`)

This is recommended for development as it automatically rebuilds and restarts the server when code changes are detected.

*   **Install `air`:** If you don't have `air`, install it:
    ```bash
    go install github.com/cosmtrek/air@latest
    ```
*   **Run `air`:**
    From the `server/` directory:
    ```bash
    air
    ```
    This will use the configuration in `.air.toml` to build and run the application. The server will typically be available at `http://localhost:5000`.

### 2. Standard Go Run

*   **Run directly:**
    From the `server/` directory:
    ```bash
    go run cmd/api/main.go
    ```
    The server will start, and you should see log output indicating it's listening on port `5000`.

## API Documentation

Detailed API endpoint descriptions, request/response formats, and usage examples can be found in the [API Specification](./api_spec.md) file.

## Testing

To run the automated tests for the backend:

1.  Navigate to the `server/` directory.
2.  Execute the following command:
    ```bash
    go test ./...
    ```
    This command will discover and run all test files (ending in `_test.go`) within the current directory and its subdirectories. The `server/internal/utils/crypto_test.go` is an example of such a test.

---
This README provides a guide to getting the Go backend up and running. For frontend setup, please refer to `angular/README.md`.
