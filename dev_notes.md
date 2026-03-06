# Development Notes

This document provides notes and instructions for setting up and running the development environment.

## Go Backend (`server` directory)

These instructions assume you are working within the `server` directory of the project.

### Installing Dependencies

To install or update the Go application dependencies, navigate to the `server` directory and run:

```bash
cd server

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

**Running with Air:**

Once `air` is installed, navigate to the `server` directory and run:

```bash
cd server
air
```

This will start the Go server, and `air` will monitor for file changes, automatically rebuilding and restarting the server as needed. Configuration is in `.air.toml`.

### Environment Configuration

Copy the example environment file and fill in values:

```bash
cd server
cp .env.example .env
```

Required environment variables include:
- `APP_SECRET`, `JWT_SECRET`, `REFRESH_SECRET`, `RESET_SECRET` - for JWT tokens
- `SMTP` settings for email (`EMAIL_HOST`, `EMAIL_PORT`, `EMAIL_USER`, `EMAIL_PASSWORD`)
- `STRIPE_API_KEY`, `STRIPE_WEBHOOK_SECRET` - for Stripe billing

## Frontend (Angular - `angular` directory)

The Angular frontend is located in the `angular/` directory.

### Setup and Running

```bash
cd angular
npm install
ng serve                 # Dev server on :4200 (proxies API to :5000)
ng build --configuration production
```

Refer to `angular/README.md` and `angular/DESIGN_SYSTEM.md` for more details.

## Codespaces Development

Frontend URL when running in GitHub Codespaces:
The specific URL can vary. Check your Codespaces port forwarding. A common pattern is:
`https://<your-codespace-name>-4200.githubpreview.dev/`

Example:
`https://juparave-go-angular-security-pgvrxv62rvj5-4200.githubpreview.dev/`
