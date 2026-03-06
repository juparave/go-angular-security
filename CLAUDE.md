# CLAUDE.md — Developer Guide

## Project Overview

`go-angular-security` is a production-ready SaaS starter template combining a **Go Fiber** backend with an **Angular 19** frontend. It ships with multi-tenancy, JWT auth, team management, Stripe subscriptions, and a transactional email system.

---

## Architecture

### Multi-tenancy
- Each `Account` has its own SQLite database in `server/data/<account_id>.db`
- A master database (`gorm.db`) stores accounts, users, and Stripe metadata
- `database.Manager` (singleton) manages connections to all tenant databases
- `UserSyncService` keeps users in sync across master and tenant DBs

### JWT Authentication
- Access tokens (short-lived) and refresh tokens (long-lived) signed with separate secrets
- Reset tokens use a third independent secret
- Configured entirely via env vars — no hardcoded secrets
- `utils.InitJWT()` must be called at startup before any JWT operations

### Stripe Billing
- Checkout sessions created via Stripe API (`/subscriptions/create-checkout-session`)
- Webhooks received at `/stripe/webhook` and `/api/v1/webhooks/stripe`
- Subscription cancel, reactivate, and plan-change endpoints in `subscription.go`

---

## Directory Structure

```
go-angular-security/
├── angular/               # Angular 19 frontend
│   ├── src/app/
│   ├── angular.json
│   └── DESIGN_SYSTEM.md   # UI patterns guide
├── server/                # Go Fiber backend
│   ├── cmd/api/main.go    # Entry point
│   ├── internal/
│   │   ├── config/        # AppConfig, JWTConfig, EmailConfig
│   │   ├── database/      # Multi-tenant DB manager
│   │   ├── handlers/      # HTTP handlers (auth, team, subscription, contact, google)
│   │   ├── middleware/    # IsAuthenticated, RequireAccount, RequireAdmin/Editor
│   │   ├── models/        # GORM models (User, Account, Role, Subscription)
│   │   ├── routes/        # Route definitions with rate limiters
│   │   ├── services/      # UserSyncService
│   │   └── utils/         # JWT, crypto, helpers
│   ├── pkg/mail/          # Async email service with templates
│   │   ├── templates/     # HTML email templates
│   │   ├── service.go     # Queue-based async sender with retry
│   │   ├── cache.go       # TTL render cache
│   │   └── smtp_sender.go # SMTP implementation
│   └── .env.example
├── Dockerfile             # Multi-stage build
├── Makefile               # Docker operations
└── CLAUDE.md              # This file
```

---

## Development Commands

### Backend
```bash
cd server
cp .env.example .env     # Fill in values
go run cmd/api/main.go   # Start dev server on :5000
go build ./...           # Verify build
go test ./...            # Run tests
```

### Frontend
```bash
cd angular
npm install
ng serve                 # Dev server on :4200 (proxies API to :5000)
ng build --configuration production
```

### Docker
```bash
make build               # Build image
make run                 # Start container (uses server/.env)
make logs                # Tail logs
make replace             # Rebuild + restart
PORT=8080 make run       # Custom port
```

---

## API Endpoints Reference

All API endpoints are prefixed `/api/v1/`.

### Public (rate limited)
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

### Protected (requires JWT)
| Method | Path | Description |
|--------|------|-------------|
| GET | `/user` | Get current user |
| POST | `/logout` | Clear auth cookies |
| PUT | `/users/info` | Update profile |
| PUT | `/users/password` | Update password |
| POST | `/change-password` | Change password (requires current) |

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

---

## Database Models

- **Account** — tenant entity with Stripe customer ID, plan tier, user limits
- **User** — belongs to Account; has role, password digest, enabled flag
- **Role** — 1=Admin, 2=Editor, 3=Viewer
- **Subscription** — mirrors Stripe subscription data

---

## Authentication Flow

1. Client sends `POST /login` with email + password
2. Server verifies, generates access token (JWT) + refresh token
3. Tokens set as `HttpOnly` cookies (`jwt` + `refreshjwt`)
4. Protected routes read token from `Token` header or `jwt` cookie
5. On expiry, client calls `POST /refresh-token` with the refresh token
6. `Secure: true` on cookies when `APP_MODE=production`

---

## Subscription Flow

1. Client calls `POST /subscriptions/create-checkout-session` with `priceId`
2. Server creates/reuses Stripe customer, returns `sessionId`
3. Client redirects to Stripe Checkout
4. On success, Stripe fires webhook → `PostStripeWebhook` updates account
5. Account tier updated based on webhook event

---

## Email System

The `pkg/mail` service provides:
- **Async queue** — emails sent in background via channel (3 workers)
- **Retry logic** — 3 attempts with 1s/2s/4s exponential backoff
- **Dev mode** — set `DevMode=true` + `TestRecipient` to redirect all emails
- **Template rendering** — HTML templates embedded at compile time

Templates: `password_reset`, `team_invitation`, `welcome`, `contact`

To add a template: create `server/pkg/mail/templates/<name>.html`.

---

## Coding Conventions

- Handlers receive `*fiber.Ctx`, return `error`
- Access `app` (config) and `emailService` via package-level vars in `repo.go`
- Multi-tenant: use `database.Manager.GetMasterDB()` for users/accounts, `database.Manager.GetConnection(accountID)` for tenant data
- Rate limiting applied at route registration in `routes.go`
- Secrets must come from env vars — never hardcode

---

## Common Issues

**"JWT not initialized"** — `utils.InitJWT()` not called before handlers execute. Check `main.go`.

**Email not sent** — `EMAIL_HOST` not set → email service is nil → emails log only. Set SMTP config.

**Tenant DB not found** — `database.Manager.GetConnection(accountID)` returns error if the tenant DB hasn't been initialized. Call `SyncUserToTenant` first.

**CORS errors** — Check `APP_DOMAIN` matches the frontend URL. In development, `localhost:4200` and `localhost:5000` are always allowed.
