# Specification: Stripe Subscriptions Integration

**Version:** 1.0
**Date:** 2024-07-27

## 1. Overview

This document outlines the requirements and implementation steps for integrating Stripe Subscriptions into an application with an Angular v18 frontend and a Go Fiber/GORM backend. The goal is to allow users to subscribe to paid plans offered by the application, manage their subscription status, and ensure the application's backend accurately reflects the subscription state stored in Stripe.

## 2. Specification

### 2.1. Functional Requirements

*   **User Subscription:** Authenticated users must be able to initiate a subscription process for a predefined Stripe Price/Plan via the Angular frontend.
*   **Stripe Checkout:** The application will utilize Stripe Checkout (hosted page) for collecting payment details and creating the subscription.
*   **Backend Handling:** The Go backend will handle:
    *   Creating/retrieving Stripe Customer objects.
    *   Creating Stripe Checkout Sessions in `subscription` mode.
    *   Storing and managing the relationship between the application's user ID and the Stripe Customer ID and Subscription ID.
    *   Securely handling Stripe API keys.
    *   Processing Stripe webhooks to keep the local database synchronized with Stripe events (subscription creation, updates, cancellations, payment failures, etc.).
*   **Database Persistence:** Subscription details (status, plan, current period end, etc.) associated with a user will be stored persistently in the application's database using GORM and a dedicated `subscriptions` table.
*   **Frontend Display:** The Angular frontend should reflect the user's current subscription status (e.g., active, inactive, canceled) based on data retrieved from the backend.
*   **Redirection:** Users should be redirected to appropriate pages within the Angular app upon successful subscription or cancellation via Stripe Checkout.

### 2.2. Non-Functional Requirements

*   **Security:** Stripe API keys must be stored securely (e.g., environment variables, secrets manager) and never exposed to the frontend. Webhook endpoints must verify Stripe signatures. All communication between frontend and backend should be over HTTPS.
*   **Reliability:** The system must reliably handle Stripe webhooks to ensure data consistency between Stripe and the local database. Implement retry mechanisms or queuing for webhook processing if necessary.
*   **User Experience:** The subscription process should be smooth and intuitive for the user. Clear feedback should be provided during and after the process.

## 3. Implementation Details

### 3.1. Frontend (Angular v18)

*   **Subscription Component/Service:**
    *   Display subscription options (plans/prices).
    *   Include a "Subscribe" button for the chosen plan.
    *   On button click, call a backend API endpoint (`/api/v1/stripe/create-checkout-session`) with the necessary details (e.g., Price ID, user identifier).
    *   Receive the Stripe Checkout Session ID from the backend.
    *   Use `@stripe/stripe-js` library's `redirectToCheckout` method with the received Session ID to redirect the user to the Stripe-hosted checkout page.
*   **Routing:**
    *   Define routes for success (`/subscription-success`) and cancellation (`/subscription-cancel`) pages. These URLs will be configured in the backend when creating the Checkout Session.
*   **Success Page (`/subscription-success`):**
    *   This page confirms to the user that the checkout process was initiated successfully.
    *   **Important:** It should *not* be the primary mechanism for confirming the subscription is active. That confirmation comes via backend webhooks.
    *   It can display a pending/processing message and potentially trigger a fetch of the user's latest subscription status from the backend after a short delay or based on user action.
*   **Cancel Page (`/subscription-cancel`):**
    *   Inform the user that the subscription process was canceled.
*   **Subscription Status Display:**
    *   Implement logic (e.g., in a user profile section or header) to fetch the user's subscription status from a dedicated backend endpoint (e.g., `/api/v1/user/status`) and display it accordingly. This data originates from the backend's `subscriptions` table.
*   **HTTP Service:** Use Angular's `HttpClient` to communicate with the backend API. Ensure authentication tokens (e.g., JWT) are included in requests.

### 3.2. Backend (Go Fiber + GORM)

*   **Configuration:**
    *   Store Stripe Secret Key and Webhook Signing Secret securely (e.g., environment variables).
    *   Configure GORM connection to the database.
*   **GORM Model (`models/subscription.go`):**
    ```go
    package models

    import (
        "time"
        "gorm.io/gorm"
    )

    type Subscription struct {
        gorm.Model // Includes ID, CreatedAt, UpdatedAt, DeletedAt

        UserID             string   `gorm:"uniqueIndex:idx_user_sub;not null"` // Foreign key to your User model
        StripeCustomerID   string `gorm:"uniqueIndex;not null"`
        StripeSubscriptionID string `gorm:"uniqueIndex:idx_user_sub;not null"`
        StripePriceID      string `gorm:"not null"`
        Status             string `gorm:"index;not null"` // e.g., "active", "canceled", "incomplete", "past_due", "trialing"
        CurrentPeriodEnd   time.Time
        TrialStart         time.Time 
        TrialEnd           time.Time 
        // Add other relevant fields as needed, e.g., CancelAtPeriodEnd bool
    }
    ```
    *   Ensure you have a corresponding `User` model with a `string ID`.
    *   Run GORM migrations to create the `subscriptions` table.
*   **Stripe Client Initialization:** Initialize the official `stripe-go` client using the secret key.
*   **API Endpoints (within Fiber router group, e.g., `/api/v1/stripe`):**
    *   **`POST /create-checkout-session`:**
        *   Requires authentication (get `userID` from context).
        *   Input: `{"priceId": "price_..."}`.
        *   Logic:
            1.  Check if the user already has an active subscription (optional, depends on business logic).
            2.  Retrieve or create a Stripe Customer:
                *   Query the `subscriptions` table (or a dedicated `user_stripe_mapping` table if preferred) for an existing `StripeCustomerID` linked to the `userID`.
                *   If not found, create a new Stripe Customer using the Stripe API (`customers.New`). Store the mapping immediately (e.g., update the `User` model or a mapping table, though the `Subscription` model will store it eventually). *Alternatively, create the customer only when the `checkout.session.completed` webhook arrives, using customer details from the session.* A common pattern is to create the customer here to associate metadata.
                *   If found, use the existing `StripeCustomerID`.
            3.  Define `success_url` (e.g., `https://yourdomain.com/subscription-success?session_id={CHECKOUT_SESSION_ID}`) and `cancel_url` (e.g., `https://yourdomain.com/subscription-cancel`).
            4.  Create a Stripe Checkout Session using `checkout/sessions.New` with:
                *   `mode: 'subscription'`
                *   `customer`: The `StripeCustomerID`.
                *   `line_items`: Containing the `priceId` from the request.
                *   `success_url`, `cancel_url`.
                *   Crucially, include `client_reference_id: userID` (or other metadata) to link the session back to your internal user upon webhook receipt.
            5.  Return `{"sessionId": "cs_..."}` to the frontend.
    *   **`POST /webhooks`:**
        *   Publicly accessible endpoint, but requires signature verification.
        *   Logic:
            1.  Read the request body.
            2.  Get the `Stripe-Signature` header.
            3.  Verify the webhook signature using `webhook.ConstructEvent` with the payload, signature header, and your webhook signing secret. Handle verification errors (return 400).
            4.  Handle specific event types (use a `switch event.Type`):
                *   **`checkout.session.completed`:**
                    *   Retrieve the `checkout.Session` object from the `event.Data.Object`.
                    *   Extract `CustomerID`, `SubscriptionID`, `client_reference_id` (your `userID`), and plan details (`PriceID`) from the session and its line items or subscription object.
                    *   Fetch the full `Subscription` object from Stripe API using the `SubscriptionID` to get the initial `Status` and `CurrentPeriodEnd`.
                    *   Create or update the record in your `subscriptions` table using GORM: Store `UserID`, `StripeCustomerID`, `StripeSubscriptionID`, `StripePriceID`, `Status`, `CurrentPeriodEnd`. Ensure atomicity (transactions if needed).
                *   **`customer.subscription.updated`:**
                    *   Retrieve the `stripe.Subscription` object.
                    *   Find the corresponding record in your `subscriptions` table using `StripeSubscriptionID`.
                    *   Update fields like `Status`, `CurrentPeriodEnd`, `StripePriceID` (if plan changed), `CancelAtPeriodEnd`, etc.
                *   **`customer.subscription.deleted`:**
                    *   Retrieve the `stripe.Subscription` object.
                    *   Find the corresponding record in your `subscriptions` table.
                    *   Update the `Status` to "canceled" (or remove the record, depending on your data retention policy).
                *   *(Optional)* Handle other events like `invoice.payment_failed`, `invoice.paid`, `customer.subscription.trial_will_end` as needed.
            5.  Return a `200 OK` response to Stripe to acknowledge receipt. If processing fails, return an appropriate error (Stripe will retry).
    *   **`GET /user/status` (Example endpoint, could be part of a general user profile endpoint):**
        *   Requires authentication (get `userID` from context).
        *   Query the `subscriptions` table for the given `userID`.
        *   Return relevant subscription details (e.g., `{"status": "active", "plan": "price_...", "expires": "..."}`) or an indication of no active subscription.

### 3.3. Data Flow (Revised based on GORM)

1.  **FRONTEND:** User clicks "Subscribe" button for a specific plan (`priceId`).
2.  **FRONTEND:** Calls backend `POST /api/v1/stripe/create-checkout-session` with `{priceId: "..."}`.
3.  **BACKEND:** (`/create-checkout-session`)
    *   Authenticates user, gets `userID`.
    *   Retrieves or creates Stripe Customer ID associated with `userID`.
    *   Creates Stripe Checkout Session (`mode: 'subscription'`, `success_url`, `cancel_url`, `line_items` with `priceId`, `customer`, `client_reference_id: userID`).
    *   Returns `{sessionId: "cs_..."}` to frontend.
4.  **FRONTEND:** Uses `@stripe/stripe-js`'s `redirectToCheckout({ sessionId: "cs_..." })`.
5.  **USER:** Completes payment/subscription on Stripe's hosted Checkout page.
6.  **STRIPE:** Redirects user to `success_url` (e.g., `/subscription-success`) or `cancel_url`.
7.  **STRIPE:** Sends `checkout.session.completed` webhook event to backend `POST /api/v1/stripe/webhooks`.
8.  **BACKEND:** (`/webhooks`)
    *   Verifies webhook signature.
    *   Parses `checkout.session.completed` event.
    *   Extracts `SubscriptionID`, `CustomerID`, `UserID` (from `client_reference_id`), `PriceID`.
    *   Fetches full subscription details from Stripe API.
    *   Creates/Updates the `Subscription` record in the database using GORM, linking `UserID` to the Stripe details and setting the initial status (`active`, `trialing`, etc.).
9.  **FRONTEND:** (`/subscription-success` page)
    *   Displays a success message.
    *   (Optional/Delayed) User navigates away or app fetches updated user status from backend (`GET /api/v1/user/status`).
10. **BACKEND:** (`/user/status`)
    *   Queries the `subscriptions` table for the authenticated user.
    *   Returns the current subscription status stored in the database.
11. **FRONTEND:** Displays the accurate subscription status based on the backend data.
12. **STRIPE:** Sends subsequent webhook events (e.g., `customer.subscription.updated`, `customer.subscription.deleted`) for renewals, cancellations, etc.
13. **BACKEND:** (`/webhooks`) Handles these events, updating the corresponding `Subscription` record in the database via GORM.

## 4. Tools & Libraries

*   **Frontend:**
    *   Angular v18+ (`@angular/core`, `@angular/common`, `@angular/router`, `@angular/common/http`, `@angular/material`, `@angular/forms`)
    *   Stripe.js (`@stripe/stripe-js`)
    *   Node.js / npm or yarn
*   **Backend:**
    *   Go (latest stable version)
    *   Fiber (`github.com/gofiber/fiber/v2`)
    *   JWT (`github.com/golang-jwt/jwt/v4`)
    *   GORM (`gorm.io/gorm` and appropriate DB driver, e.g., `gorm.io/driver/sqlite`)
    *   Stripe Go library (`github.com/stripe/stripe-go/v8x` - replace `x` with latest version, e.g., v81)
    *   Environment variable library (e.g., `github.com/joho/godotenv`) for sensitive key management
    *   Load from `config.json` using `github.com/golobby/config/v3` for configuration management
*   **Database:** PostgreSQL (recommended), MySQL, or SQLite (for development)
*   **Stripe:**
    *   Stripe Account
    *   Stripe API Keys (Publishable Key for Frontend, Secret Key for Backend)
    *   Stripe Product and Price IDs configured in the Stripe Dashboard.
    *   Stripe CLI (for local webhook testing)
*   **Development:**
    *   IDE (e.g., VS Code with Go and Angular extensions)
    *   Git / Version Control

## 5. Step-by-Step Procedure

1.  **Setup Stripe:**
    *   Create a Stripe account.
    *   Define Products and corresponding Prices (Subscription type) in the Stripe Dashboard. Note the Price IDs (`price_...`).
    *   Obtain your Publishable Key, Secret Key, and create a Webhook Endpoint Signing Secret.
2.  **Backend Setup:**
    *   Initialize Go project with Fiber.
    *   Install dependencies: `go get github.com/gofiber/fiber/v2 gorm.io/gorm gorm.io/driver/sqlite github.com/stripe/stripe-go/v81 github.com/joho/godotenv`
    *   Configure environment variables (`.env` file) for `STRIPE_SECRET_KEY`, `STRIPE_WEBHOOK_SECRET`, `DATABASE_URL`.
    *   Set up database connection and GORM initialization.
    *   Define `User` and `Subscription` GORM models. Run auto-migrations.
    *   Implement user authentication middleware (if not already present).
3.  **Backend Implementation:**
    *   Initialize Stripe client in a suitable package/service.
    *   Implement the `POST /api/v1/stripe/create-checkout-session` handler.
    *   Implement the `POST /api/v1/stripe/webhooks` handler, including signature verification and event handling logic (`checkout.session.completed`, `customer.subscription.updated`, `customer.subscription.deleted`).
    *   Implement the `GET /api/v1/user/status` (or similar) endpoint to fetch subscription data from the DB.
4.  **Frontend Setup:**
    *   Initialize Angular project (if new) or navigate to existing project.
    *   Install dependencies: `npm install @stripe/stripe-js`.
    *   Configure environment variables (e.g., `environment.ts`) for `STRIPE_PUBLISHABLE_KEY`.
5.  **Frontend Implementation:**
    *   Create an `HttpClient` service for backend communication.
    *   Create a component to display subscription options and the "Subscribe" button.
    *   Implement the button's click handler:
        *   Call the backend `/create-checkout-session` endpoint.
        *   Load Stripe.js using `loadStripe`.
        *   Call `stripe.redirectToCheckout({ sessionId })`.
    *   Create routing and components for `/subscription-success` and `/subscription-cancel` pages.
    *   Implement logic to display the user's subscription status by fetching data from the backend `/user/status` endpoint.
6.  **Webhook Configuration:**
    *   Start your backend server.
    *   Use Stripe CLI to forward webhooks to your local development server: `stripe listen --forward-to localhost:PORT/api/v1/stripe/webhooks` (replace `PORT` with your backend port).
    *   Use the webhook signing secret provided by the Stripe CLI command for local testing.
    *   In your Stripe Dashboard (Test mode first), create a webhook endpoint pointing to your *deployed* backend URL (`https://yourdomain.com/api/v1/stripe/webhooks`). Select the necessary events (at least `checkout.session.completed`, `customer.subscription.updated`, `customer.subscription.deleted`). Add the production webhook signing secret to your deployed backend's environment variables.
7.  **Testing:**
    *   Use Stripe Test Cards to simulate the checkout flow.
    *   Verify that the `checkout.session.completed` webhook is received, processed correctly, and the `subscriptions` table is populated.
    *   Test subscription updates and cancellations using the Stripe Dashboard (Test mode) and verify the corresponding webhooks update the database correctly.
    *   Verify the frontend accurately reflects the subscription status fetched from the backend.
    *   Test error handling (e.g., invalid API keys, webhook signature failure).
8.  **Deployment:**
    *   Deploy the backend and frontend.
    *   Ensure all necessary environment variables (Stripe keys, DB connection, webhook secret) are configured correctly in the production environment.
    *   Configure the production webhook endpoint in the Stripe Dashboard.

## 6. Security Considerations

*   **Never** expose your Stripe Secret Key or Webhook Signing Secret in frontend code or insecurely on the backend.
*   **Always** verify webhook signatures on the backend before processing the event.
*   Use HTTPS for all communication.
*   Protect backend endpoints with authentication and authorization, ensuring users can only affect their own data.
*   Sanitize any user input used in API calls.
