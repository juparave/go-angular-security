import { inject } from '@angular/core';
import {
  ActivatedRouteSnapshot,
  RouterStateSnapshot,
  Router,
  CanActivateFn,
} from '@angular/router';
import { Store, select } from '@ngrx/store';
import { map, Observable, tap, filter, take } from 'rxjs';
import { AppState } from 'src/app/store/interfaces/app-state';
import { selectCurrentUser } from 'src/app/store/selectors/auth.selectors';
import { User } from 'src/app/model/user';

// --- Base Logic ---

/**
 * Base function to create a subscription status guard.
 * Checks if the current user's subscription status is one of the allowed statuses.
 * Redirects to '/app/subscription' if the status is not allowed or user/subscription is missing.
 *
 * @param allowedStatuses Array of allowed subscription statuses.
 * @returns CanActivateFn
 */
const createSubscriptionGuard = (allowedStatuses: string[]): CanActivateFn => {
  return (route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> => {
    const store = inject(Store<AppState>);
    const router = inject(Router);

    return store.pipe(
      select(selectCurrentUser),
      take(1), // Ensure the observable completes after getting the first value
      map((user: User | null) => {
        const currentStatus = user?.subscription?.status;
        const isAllowed = !!currentStatus && allowedStatuses.includes(currentStatus);

        if (!isAllowed) {
          // Redirect if status is not allowed, or if user/subscription is missing
          // Consider redirecting to a specific subscription page if available
          router.navigate(['/app/subscription']); // Or '/login' or another appropriate route
          return false;
        }
        return true;
      })
    );
  };
};

// --- Specific Guards ---

/**
 * Guard: Allows access only if the user's subscription status is 'active' or 'canceled'.
 */
export const canActivateActiveCanceledSubscription: CanActivateFn = createSubscriptionGuard(['active', 'canceled']);

/**
 * Guard: Allows access only if the user's subscription status is 'active', 'trialing', or 'canceled'.
 */
export const canActivateActiveTrialingCanceledSubscription: CanActivateFn = createSubscriptionGuard(['active', 'trialing', 'canceled']);
