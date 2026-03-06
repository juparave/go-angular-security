# Go Angular Security

A production-ready SaaS template built with Go backend and Angular frontend, featuring multitenancy, subscription billing, team management, and modern UI.

## Features

- **Multitenancy** - SQLite-based tenant isolation with connection pooling
- **Modern UI** - Angular 19 + TailwindCSS v4 + DaisyUI
- **Authentication** - JWT-based auth with password reset and email verification
- **Team Management** - Role-based access control (Admin/Editor/Viewer)
- **Subscription Billing** - Stripe integration with 2-tier Free/Pro plans
- **Email System** - Async email service with templates

## Prerequisites

- Go 1.21 or higher
- Node.js 20+ and npm
- Angular CLI 19+
- Stripe account (for payments)

## Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/juparave/go-angular-security.git
cd go-angular-security
```

### 2. Backend Setup

```bash
cd server
cp .env.example .env
# Edit .env with your configuration
go mod tidy
go run cmd/api/main.go
```

The Go backend will start on `http://localhost:5000`.

### 3. Frontend Setup

```bash
cd angular
npm install
ng serve
```

The Angular frontend will start on `http://localhost:4200`.

## Architecture

### Multitenancy

Each tenant (account) has its own SQLite database stored at `data/{accountID}/data.db`. The master database stores account information and user credentials for authentication.

```
data/
├── master.db          # Account and user registry
├── acc_abc123/        # Tenant database
│   └── data.db
└── acc_xyz789/
    └── data.db
```

### Subscription Tiers

| Feature | Free | Pro |
|---------|------|-----|
| Users | 1 | 10 |
| Storage | 100MB | 10GB |
| Support | Community | Priority |

### Roles & Permissions

- **Admin** - Full access, team management, billing
- **Editor** - Create, edit, and view content
- **Viewer** - Read-only access

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_NAME` | Application name | GO-ANGULAR-SECURITY |
| `APP_PORT` | Server port | 5000 |
| `APP_DOMAIN` | Frontend URL | http://localhost:4200 |
| `JWT_SECRET` | Access token secret | (required) |
| `JWT_REFRESH_SECRET` | Refresh token secret | (required) |
| `DATABASE_PATH` | Master DB path | gorm.db |
| `EMAIL_HOST` | SMTP host | - |
| `EMAIL_PORT` | SMTP port | 587 |
| `STRIPE_SECRET_KEY` | Stripe API key | (required) |

See `.env.example` for all options.

## API Endpoints

### Public

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/register-account` | Create new account |
| POST | `/api/v1/login` | Authenticate user |
| POST | `/api/v1/request-password-reset` | Request reset email |
| POST | `/api/v1/reset-password` | Reset with token |

### Authenticated

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/user` | Current user |
| POST | `/api/v1/logout` | Sign out |
| POST | `/api/v1/change-password` | Change password |

### Team Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/team` | List members |
| POST | `/api/v1/team` | Invite member |
| PUT | `/api/v1/team/:id` | Update member |
| DELETE | `/api/v1/team/:id` | Remove member |

### Subscriptions

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/subscriptions/current` | Current subscription |
| POST | `/api/v1/subscriptions/create-checkout-session` | Stripe checkout |

## Development

### Backend

```bash
cd server
# Run tests
go test ./...

# Run with hot reload (requires air)
air
```

### Frontend

```bash
cd angular
# Development server
ng serve

# Build for production
ng build

# Run tests
ng test
```

## Deployment

### Docker

```bash
# Build
docker-compose build

# Run
docker-compose up -d
```

### Manual

1. Build the Go binary: `go build -o server ./cmd/api`
2. Build Angular: `ng build --configuration production`
3. Configure environment variables
4. Run the server with the `dist/angular` folder served statically

## License

MIT License - see LICENSE file for details.
