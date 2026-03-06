# go-angular-security

A production-ready SaaS starter template — **Go Fiber** backend + **Angular 19** frontend — with multi-tenancy, JWT authentication, team management, Stripe subscriptions, and a transactional email system.

---

## Features

- **Multi-tenancy** — Per-tenant SQLite databases with connection pooling
- **JWT Auth** — Access + refresh token pair, config-driven secrets, secure cookies
- **Google OAuth** — Sign in with Google via frontend flow
- **Password Reset** — Token-based reset with digest-based reuse prevention
- **Team Management** — Invite members, assign roles (Admin / Editor / Viewer)
- **Stripe Billing** — Checkout sessions, webhooks, cancel/reactivate/plan-change
- **Email System** — Async queue, 3-attempt retry, dev-mode redirect, HTML templates
- **Rate Limiting** — Per-IP limits on auth and password reset endpoints
- **Graceful Shutdown** — SIGTERM/SIGINT handling with queue drain
- **Contact Form** — Email notification with configurable recipient
- **Docker** — Multi-stage build producing a minimal Alpine image
- **Modern UI** — Angular 19, Tailwind CSS v4, DaisyUI 5

---

## Tech Stack

| Layer | Technology | Version |
|-------|-----------|---------|
| Backend | Go + Fiber | 1.24 / v2 |
| Frontend | Angular | 19 |
| Database | SQLite (GORM) | — |
| Auth | golang-jwt | v4 |
| Billing | Stripe | v82 |
| UI | Tailwind CSS + DaisyUI | v4 / 5 |

---

## Prerequisites

- Go 1.24+
- Node.js 22+ and npm
- Angular CLI 19+
- Stripe account (optional — required for billing features)
- SMTP server (optional — required for email features)

---

## Quick Start

### 1. Clone

```bash
git clone https://github.com/juparave/go-angular-security.git
cd go-angular-security
```

### 2. Configure

```bash
cd server
cp .env.example .env
# Edit .env — set JWT_SECRET, JWT_REFRESH_SECRET, JWT_RESET_SECRET at minimum
```

### 3. Start the backend

```bash
go mod tidy
go run cmd/api/main.go
# → listening on http://localhost:5000
```

### 4. Start the frontend

```bash
cd ../angular
npm install
ng serve
# → http://localhost:4200
```

---

## Environment Configuration

See `server/.env.example` for the full list. Key variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_MODE` | `development` or `production` | `development` |
| `APP_PORT` | Server port | `5000` |
| `APP_DOMAIN` | Frontend URL (used in emails, CORS) | `http://localhost:4200` |
| `JWT_SECRET` | Access token signing secret | **required** |
| `JWT_REFRESH_SECRET` | Refresh token secret | **required** |
| `JWT_RESET_SECRET` | Password reset token secret | **required** |
| `JWT_ACCESS_TOKEN_TTL` | Access token lifetime (hours) | `24` |
| `JWT_REFRESH_TOKEN_TTL` | Refresh token lifetime (hours) | `720` |
| `DATABASE_PATH` | Master SQLite path | `gorm.db` |
| `EMAIL_HOST` | SMTP host | — |
| `EMAIL_PORT` | SMTP port | `587` |
| `STRIPE_SECRET_KEY` | Stripe secret key | — |
| `STRIPE_WEBHOOK_SECRET` | Stripe webhook signing secret | — |
| `GOOGLE_CLIENT_ID` | Google OAuth client ID (frontend) | — |

---

## Docker Deployment

```bash
# Build and start
make build
make run

# Tail logs
make logs

# Rebuild and restart
make replace

# Custom port
PORT=8080 make run
```

---

## Architecture Overview

```
Client (Angular)
       │
       ▼
Go Fiber API (/api/v1/*)
       │
       ├── Master DB (gorm.db) — accounts, users, Stripe metadata
       │
       └── Tenant DBs (data/<account_id>.db) — per-account data
```

- Protected routes require a valid `jwt` cookie or `Token` header
- Cookies are `Secure: true` when `APP_MODE=production`
- CORS allows `APP_DOMAIN` + localhost origins in development

### Roles & Permissions

- **Admin** — Full access, team management, billing
- **Editor** — Create, edit, and view content
- **Viewer** — Read-only access

---

## API Reference

Full reference in `CLAUDE.md`. Quick overview:

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/register-account` | Create account + first user |
| POST | `/api/v1/login` | Email/password login |
| POST | `/api/v1/glogin` | Google sign-in |
| POST | `/api/v1/request-password-reset` | Send reset email |
| POST | `/api/v1/reset-password` | Reset with token |
| GET | `/api/v1/user` | Current user (auth required) |
| GET/POST/PUT/DELETE | `/api/v1/team` | Team management |
| GET | `/api/v1/subscriptions/current` | Active Stripe subscription |
| POST | `/api/v1/subscriptions/create-checkout-session` | Stripe checkout |
| POST | `/api/v1/contact` | Contact form |

---

## Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/my-feature`
3. Commit your changes
4. Open a pull request

---

## License

MIT — see `LICENSE` for details.
