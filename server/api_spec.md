# API Specification

## Overview

The API is built with Go using the Fiber framework. It provides endpoints for user registration, authentication, profile management, and (planned) subscription management. All endpoints are grouped under `/api/v1`.

## Endpoints

### Public Endpoints

- `POST /api/v1/register`  
  Registers a new user. Expects user details and password confirmation. Returns the created user object and sets authentication cookies.

- `POST /api/v1/login`  
  Authenticates a user with email and password. Returns the user object and sets authentication cookies.

### Authenticated Endpoints

All endpoints below require authentication (JWT in cookie or header):

- `PUT /api/v1/users/info`  
  Updates the authenticated user's profile information (first name, last name, email).

- `PUT /api/v1/users/password`  
  Updates the authenticated user's password. Requires password confirmation.

- `POST /api/v1/logout`  
  Logs out the user by expiring the authentication cookie.

- `GET /api/v1/user`  
  Returns the authenticated user's profile information.

### Planned/To Be Implemented

- `POST /api/v1/stripe/create-checkout-session`  
  Initiates a Stripe Checkout session for subscriptions.

- `POST /api/v1/stripe/webhooks`  
  Handles Stripe webhook events for subscription status updates.

- `GET /api/v1/user/status`  
  Returns the authenticated user's current subscription status.

## Authentication

- Uses JWT stored in HTTP-only cookies.
- Middleware checks for valid JWT and user existence.

## Security

- Passwords are hashed before storage.
- All sensitive operations require authentication.
- CORS is configured to allow requests from the Angular frontend.

## Static Content

- `GET /api/v1/uploads/*`  
  Serves static files from the `uploads` directory.
