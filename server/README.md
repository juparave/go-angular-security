# Go Backend API

This directory (`server/`) contains the Go backend API for the go-angular-security SaaS starter. It's built using the [Fiber](https://docs.gofiber.io/) web framework and provides RESTful endpoints for the Angular frontend.

## Project Overview

The Go backend is responsible for:

- User authentication (registration, login, Google OAuth) and JWT generation/validation
- Authorization and role-based access control (Admin / Editor / Viewer)
- Multi-tenant database management with per-account SQLite databases
- Team management and member invitations
- Stripe subscription billing and webhook handling
- Transactional email system with async queue and retry logic

## Prerequisites

- Go 1.24+
- SMTP server (optional — required for email features)
- Stripe account (optional — required for billing features)

## Setup and Running

### 1. Navigate to the server directory

```bash
cd server
```

### 2. Install Dependencies

```bash
go mod tidy
go mod download
```

### 3. Environment Configuration

Copy the example environment file and configure:

```bash
cp .env.example .env
```

Required variables (see `.env.example` for full list):

| Variable | Description |
|----------|-------------|
| `JWT_SECRET` | Access token signing secret |
| `JWT_REFRESH_SECRET` | Refresh token secret |
| `JWT_RESET_SECRET` | Password reset token secret |
| `EMAIL_HOST` | SMTP host (optional) |
| `STRIPE_SECRET_KEY` | Stripe API key (optional) |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID (optional) |

### 4. Run the Application

**Without auto-reload:**

```bash
go run cmd/api/main.go
# → listening on http://localhost:5000
```

**With auto-reload using Air (recommended for development):**

```bash
# Install Air if not already installed
go install github.com/cosmtrek/air@latest

# Run with auto-reload
air
# Uses .air.toml configuration
```

## API Endpoints

All API endpoints are prefixed `/api/v1/`. Routes are defined in `internal/routes/routes.go`.

### Public Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/register` | Legacy user registration |
| POST | `/register-account` | Self-service account + user creation |
| POST | `/login` | Email/password login |
| POST | `/glogin` | Google OAuth sign-in |
| POST | `/refresh-token` | Refresh JWT access token |
| POST | `/request-password-reset` | Send reset email |
| POST | `/reset-password` | Reset password with token |
| POST | `/contact` | Contact form submission |

### Protected Endpoints (requires JWT)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/user` | Get current user |
| POST | `/logout` | Clear auth cookies |
| PUT | `/users/info` | Update profile |
| PUT | `/users/password` | Update password |

### Team Management

| Method | Path | Description |
|--------|------|-------------|
| GET | `/team` | List team members |
| POST | `/team` | Invite new member |
| PUT | `/team/:id` | Update member role/status |
| DELETE | `/team/:id` | Remove member |
| POST | `/team/:id/resend` | Resend invitation email |

### Subscriptions

| Method | Path | Description |
|--------|------|-------------|
| GET | `/subscriptions/current` | Get active subscription |
| POST | `/subscriptions/create-checkout-session` | Create Stripe checkout |
| PATCH | `/subscriptions/:id` | Update metadata |
| POST | `/subscriptions/:id/cancel` | Cancel at period end |
| POST | `/subscriptions/:id/reactivate` | Remove cancel flag |
| POST | `/subscriptions/:id/change-plan` | Switch price/plan |

## Database Architecture

### Multi-tenancy

- **Master DB** (`gorm.db`): Stores accounts, users, and Stripe metadata
- **Tenant DBs** (`data/<account_id>.db`): Per-account SQLite databases

The `database.Manager` singleton manages connections to all tenant databases. Use:

- `database.Manager.GetMasterDB()` for users/accounts
- `database.Manager.GetConnection(accountID)` for tenant data

### Models

Key models in `internal/models/`:

- **Account** — Tenant entity with Stripe customer ID, plan tier, user limits
- **User** — Belongs to Account; has role, password digest, enabled flag
- **Role** — 1=Admin, 2=Editor, 3=Viewer
- **Subscription** — Mirrors Stripe subscription data

## Authentication

### JWT Tokens

- **Access token** — Short-lived, signed with `JWT_SECRET`
- **Refresh token** — Long-lived, signed with `JWT_REFRESH_SECRET`
- **Reset token** — For password resets, signed with `JWT_RESET_SECRET`

Tokens are set as `HttpOnly` cookies (`jwt` + `refreshjwt`). Protected routes read from `Token` header or `jwt` cookie.

### Middleware

- `IsAuthenticated` — Validates JWT token
- `RequireAccount` — Ensures user belongs to an account
- `RequireAdmin` / `RequireEditor` — Role-based access control

## Project Structure

```
server/
├── cmd/api/main.go        # Entry point
├── internal/
│   ├── config/            # AppConfig, JWTConfig, EmailConfig
│   ├── database/          # Multi-tenant DB manager
│   ├── handlers/          # HTTP handlers (auth, team, subscription, contact)
│   ├── middleware/        # Auth and role middleware
│   ├── models/            # GORM models
│   ├── routes/            # Route definitions with rate limiters
│   ├── services/          # UserSyncService
│   └── utils/             # JWT, crypto, helpers
├── pkg/mail/              # Async email service with templates
│   ├── templates/         # HTML email templates
│   ├── service.go         # Queue-based async sender with retry
│   └── smtp_sender.go     # SMTP implementation
├── .air.toml              # Air live-reload configuration
├── .env.example           # Environment variables template
└── go.mod                 # Go module definition
```

## Email System

The `pkg/mail` service provides:

- **Async queue** — Emails sent in background via channel (3 workers)
- **Retry logic** — 3 attempts with exponential backoff (1s/2s/4s)
- **Dev mode** — Set `DevMode=true` + `TestRecipient` to redirect all emails
- **Template rendering** — HTML templates embedded at compile time

Available templates: `password_reset`, `team_invitation`, `welcome`, `contact`

## Testing

```bash
go test ./...
```

## Further Reading

- Root `CLAUDE.md` — Full architecture and API reference
- Root `readme.md` — Project overview and quick start
