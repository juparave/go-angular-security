# Development Notes

This document provides notes and instructions for setting up and running the development environment.

## Go Backend (`api` directory)

These instructions assume you are working within the `api` directory of the project.

### Initial Go Module Setup (if starting from scratch)

If you were creating the Go module from scratch (this is already done in the project):

```bash
# Create a directory for the Go application
mkdir go-app
cd go-app

# Initialize the Go module
go mod init go-app
```

### Installing Dependencies

To install or update the Go application dependencies, navigate to the `api` directory and run:

```bash
# Download and install dependencies listed in go.mod
go mod download

# Tidy up the go.mod and go.sum files, removing unused dependencies
go mod tidy
```

The primary web framework used is Fiber v2. If you need to install it manually (though `go mod tidy` should handle it):
```bash
go get github.com/gofiber/fiber/v2
```

### Development Server with Auto-Reload (Air)

[Air](https://github.com/cosmtrek/air) is a tool for live-reloading Go applications. This helps in development by automatically rebuilding and restarting the server when code changes are detected.

**Installation (if not already installed):**

The `air` executable should ideally be in your system's `PATH`. You can install it globally or locally. For example, to install it in your Go bin directory (which should be in your PATH):
```bash
go install github.com/cosmtrek/air@latest
```
Alternatively, you can download a pre-compiled binary from the [Air releases page](https://github.com/cosmtrek/air/releases) and place it in a directory included in your `PATH` (e.g., `~/bin/` or `/usr/local/bin/`).

**Running with Air:**

Once `air` is installed, navigate to the `api` directory and run:

```bash
air
```

This will start the Go server, and `air` will monitor for file changes, automatically rebuilding and restarting the server as needed. Configuration for `air` can be done via a `.air.toml` file (if one exists in the directory).

## Frontend (Angular - `frontend` directory)

Refer to the main `readme.md` for instructions on setting up and running the Angular frontend.

## Codespaces Development

Frontend URL when running in GitHub Codespaces:
The specific URL can vary. Check your Codespaces port forwarding. A common pattern is:
`https://<your-codespace-name>-4200.githubpreview.dev/`

Example:
`https://juparave-go-angular-security-pgvrxv62rvj5-4200.githubpreview.dev/`
