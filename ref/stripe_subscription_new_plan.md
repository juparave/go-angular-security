# Plan for Implementing Subscription Page Flow

**1. Create the Subscription Page:**

*   **Generate Component/Module:** Create a new Angular module (e.g., `SubscriptionModule`) and a component (e.g., `SubscriptionPageComponent`) within it. This page will be responsible for displaying subscription options and handling the subscription process (likely interacting with `SubscriptionService`).
*   **Define Route:** Add a new route, say `/subscription`, in `angular/src/app/app-routing.module.ts`. This route will load the `SubscriptionModule`. It needs to be protected by the `AuthGuard` to ensure only logged-in users can access it.

**2. Update Routing Logic & Guards:**

*   **Modify `/app` Route:** In `angular/src/app/app-routing.module.ts`, update the existing `/app` route. It currently uses `canActivate: [AuthGuard]`. Change this to `canActivate: [AuthGuard, canActivateAnySubscription]`. This ensures that users must be logged in *and* have an active subscription (checked by `canActivateAnySubscription`) to access the main application.
*   **Update `SubscriptionGuard` Redirect:** Modify the `canActivateAnySubscription` function within `angular/src/app/services/auth/subscription.guard.ts`. Change the redirect path from `/renew-subscription` (or `/upgrade`) to the new `/subscription` route.

**3. Synchronize NgRx State:**

*   **Problem:** The `SubscriptionGuard` relies on the `SubscriptionState`, but the subscription data is initially fetched as part of the `User` object during login and stored in the `AuthState`. These two states are currently separate.
*   **Solution:** Create a new NgRx effect (e.g., `SyncSubscriptionEffect`).
    *   This effect will listen for `loginSuccessAction` and `getCurrentUserSuccessAction` (from `auth.actions.ts`).
    *   When these actions occur, the effect will extract the `subscription` object from the `action.currentUser` payload.
    *   It will then dispatch an existing or new action (e.g., `getSubscriptionSuccessAction` or a new `setSubscriptionFromUserAction`) with the extracted subscription data as the payload.
    *   The `subscriptionReducer` (`subscription.reducers.ts`) needs to handle this action to update the `SubscriptionState`. This ensures the `SubscriptionState` is populated correctly right after login or user refresh, making the data available for the `SubscriptionGuard`.

**4. Implement Subscription Page Functionality:**

*   Develop the UI for the `SubscriptionPageComponent` to display available plans.
*   Use the existing `SubscriptionService` (`angular/src/app/services/subscription.service.ts`) and NgRx subscription actions/effects (`angular/src/app/store/actions/subscription.actions.ts`, `angular/src/app/store/effects/subscription.effects.ts`) to handle:
    *   Fetching available subscription plans from the backend.
    *   Initiating the subscription process (e.g., redirecting to Stripe Checkout via the backend).
    *   Handling successful subscription updates (which should eventually allow the `canActivateAnySubscription` guard to pass).

**5. Backend Considerations (Go Fiber/Gorm):**

*   **Verify Login Response:** Double-check that the Go backend's login endpoint (`/internal/handlers/auth.go`) correctly includes the user's `Subscription` details within the returned `User` object.
*   **Subscription Endpoints:** Ensure the necessary backend endpoints exist and function correctly:
    *   Fetching subscription plans.
    *   Creating subscription checkout sessions (e.g., Stripe).
    *   Handling payment provider webhooks (`/internal/handlers/webhooks.go`) to update subscription status (`active`, `canceled`, etc.) in the database (`/internal/models/subscription.go`).

**Flow Summary:**

1.  User logs in.
2.  `LoginEffect` fetches `User` (with `Subscription`) -> `loginSuccessAction` dispatched.
3.  `AuthReducer` updates `AuthState`.
4.  *New* `SyncSubscriptionEffect` listens to `loginSuccessAction`, extracts `subscription`, dispatches `getSubscriptionSuccessAction` (or similar).
5.  `SubscriptionReducer` updates `SubscriptionState`.
6.  Login effect redirects to the intended URL (e.g., `/app`).
7.  Router attempts to activate `/app`.
8.  `AuthGuard` runs (passes, user is logged in).
9.  `canActivateAnySubscription` guard runs:
    *   **If Active Subscription:** Reads `SubscriptionState` (now populated), finds an active subscription, allows access to `/app`.
    *   **If No Active Subscription:** Reads `SubscriptionState`, finds no active subscription, redirects to `/subscription`.
10. User interacts with the `/subscription` page to choose and activate a plan.
