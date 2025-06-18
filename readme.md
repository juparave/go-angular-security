# Go and Angular with Security

This project demonstrates a secure web application built with a Go backend and an Angular frontend. It includes features like user authentication and authorization, CSRF protection, and secure communication between the frontend and backend.

## Prerequisites

Before you begin, ensure you have the following installed:

* Go (version 1.19 or higher)
* Node.js and npm (Node version 16 or higher, npm version 8 or higher)
* Angular CLI (version 15 or higher)

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/juparave/go-angular-security.git
   cd go-angular-security
   ```

2. **Backend Setup:**
   ```bash
   cd api
   go mod tidy
   go run main.go
   ```
   The Go backend will start on `http://localhost:3000`.

3. **Frontend Setup:**
   ```bash
   cd ../frontend
   npm install
   ng serve
   ```
   The Angular frontend will start on `http://localhost:4200`.

## Usage

Once both the backend and frontend are running, you can access the application in your web browser at `http://localhost:4200`.

The application will demonstrate secure user registration, login, and access to protected resources.

## Frontend URL for Codespaces

If you are running this project in a GitHub Codespace, the frontend might be accessible via a different URL. Please check your Codespace's port forwarding settings. A common pattern for the URL is:
`https://<your-codespace-name>-4200.githubpreview.dev/`

For example:
`https://juparave-go-angular-security-pgvrxv62rvj5-4200.githubpreview.dev/`

## Go API

The backend is built using [Fiber](https://docs.gofiber.io/), a Go web framework inspired by Express.js.

## VSCode Settings (Optional)

If you are using VSCode with the Go extension, you might find the following setting useful in your `settings.json` for multi-module workspace support:

```json
"gopls": {
    "experimentalWorkspaceModule": true
}
```
