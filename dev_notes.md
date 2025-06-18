# Development Notes

This document contains general notes and tips for developing this project. For specific setup and usage instructions for each part of the application, please refer to their respective README files:

*   **Main Project Overview:** See `readme.md`
*   **Angular Frontend:** See `angular/README.md`
*   **Go Backend Server:** See `server/README.md`

## Development Environment

### GitHub Codespaces

*   **Frontend URL:** When running the Angular development server in Codespaces, it's typically accessible via a URL like: `https://[YOUR_CODESPACE_NAME]-4200.githubpreview.dev/`
    *   The example URL previously noted was: `https://juparave-go-angular-security-pgvrxv62rvj5-4200.githubpreview.dev/` (This will vary per Codespace instance).
*   **Backend API:** The Go backend server will run on a different port (e.g., `5000` as configured) within the Codespace.

### General Workflow

*   Ensure both the Angular frontend development server and the Go backend server are running simultaneously for a complete development experience.
*   Refer to the `server/README.md` for instructions on running the Go backend, including using `air` for live reloading.
*   Refer to the `angular/README.md` for instructions on running the Angular frontend using `ng serve`.

---
*Previous notes regarding specific Go module initialization, dependency installation, and `air` setup have been moved to `server/README.md` for better organization.*
