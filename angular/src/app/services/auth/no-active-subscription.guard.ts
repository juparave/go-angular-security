import { inject } from '@angular/core';
import { CanActivateFn, Router, UrlTree, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Observable, map, filter, combineLatest } from 'rxjs';
import { Store } from '@ngrx/store';
import { AppState } from 'src/app/store/interfaces/app-state';
import { selectIsSubscriptionActive } from 'src/app/store/selectors/subscription.selectors';
import { selectIsLoggedIn } from 'src/app/store/selectors/auth.selectors';

/**
 * Guard to prevent users with an active subscription from accessing certain routes (e.g., /subscription).
 * Ensures both auth and subscription states are loaded before making a decision.
 */
export const canActivateNoActiveSubscription: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
): Observable<boolean | UrlTree> => {
  const store = inject(Store<AppState>);
  const router = inject(Router);

  return combineLatest([
    store.select(selectIsLoggedIn),
    store.select(selectIsSubscriptionActive)
  ]).pipe(
    // Wait until both auth and subscription states are no longer loading
    filter(([authState, subState]) => !authState.isLoading && !subState.isLoading),
    map(([authState, subState]) => {
      if (!authState.isLoggedIn) {
        // User is not logged in, redirect to login
        return router.createUrlTree(['/login']);
      } else if (subState.isActive) {
        // User is logged in AND has an active subscription, redirect away from subscription page
        return router.createUrlTree(['/app']);
      } else {
        // User is logged in and does NOT have an active subscription, allow access
        return true;
      }
    })
  );
};
