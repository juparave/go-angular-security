# Go and Angular with Security

This project demonstrates a full-stack web application featuring a Go backend and an Angular frontend, with a strong emphasis on security aspects like user authentication and role-based access control.

## Architecture

*   **Frontend:** Developed with Angular, providing a dynamic and responsive user interface.
*   **Backend:** Built with Go using the [Fiber](https://docs.gofiber.io/) framework, offering a high-performance API.

## Key Features

*   Secure user authentication
*   Role-based access control (RBAC)
*   JWT-based session management
*   Password hashing

## Setup and Running

### Frontend (Angular)

1.  Navigate to the `angular` directory: `cd angular`
2.  Install dependencies: `npm install` (if not already done)
3.  Run the development server: `ng serve`
4.  Open your browser and navigate to `http://localhost:4200/`.

For more detailed Angular CLI commands and information, please refer to `angular/README.md`.

### Backend (Go)

1.  Navigate to the `server` directory: `cd server`
2.  Ensure you have Go installed on your system.
3.  Install backend dependencies: `go mod tidy`
4.  You can run the application using: `go run cmd/api/main.go`
    *   This will typically start the server on port `5000` (as configured in `server/cmd/api/main.go`, check console output for any changes).
5.  Alternatively, you can build the application: `go build -o main cmd/api/main.go` and then run the executable `./main`.

For comprehensive backend setup instructions, running the server, project structure, and more, please refer to `server/README.md`. Detailed API specifications can be found in `server/api_spec.md`.

## Codespaces

If you are using GitHub Codespaces:

*   **Frontend URL:** The Angular development server will typically be available at a URL provided by Codespaces (often ending in `-4200.githubpreview.dev/`). Check the "Ports" tab in your Codespace for the exact URL.
*   **Backend API:** The Go API will be running on another port (e.g., 8080).

## VSCode Settings

For optimal Go development experience in VSCode, especially when working with modules, add the following to your `.vscode/settings.json` file (create the file/directory if it doesn't exist):

```json
{
    "gopls": {
        "experimentalWorkspaceModule": true
    }
}
```

This project aims to provide a solid foundation for building secure web applications. We welcome contributions and feedback from the community!
---

*Note: The previous content regarding "Go Api" using Fiber is now integrated into the "Backend (Go)" and "Architecture" sections.*
