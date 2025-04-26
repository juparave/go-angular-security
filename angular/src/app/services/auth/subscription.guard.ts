import { inject, Injectable } from '@angular/core';
import { CanActivateFn, Router, UrlTree, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Observable, map, filter } from 'rxjs';
import { Store } from '@ngrx/store';
import { AppState } from 'src/app/store/interfaces/app-state';
import {
  selectHasActiveTrial,
  selectHasActiveBasic,
  selectHasActivePro,
  selectHasActivePaid,
  selectIsSubscriptionActive
} from 'src/app/store/selectors/subscription.selectors';

/**
 * Base Subscription Guard class that provides subscription verification functionality
 */
@Injectable({
  providedIn: 'root',
})
export class SubscriptionGuard {
  constructor(private store: Store<AppState>, public router: Router) { }

  /**
   * Creates a guard function that verifies subscription requirements
   * @param selector The selector to use for subscription verification
   * @param redirectTo The route to redirect to if the subscription check fails
   */
  createGuard(
    selector: any,
    redirectTo: string = '/upgrade'
  ) {
    return (
      route: ActivatedRouteSnapshot,
      state: RouterStateSnapshot
    ): Observable<boolean | UrlTree> => {
      return this.store.select(selector).pipe(
        filter(state => !state.isLoading),
        map(result => {
          // Extract the boolean result from the selector result
          // Each selector returns an object with a different key that indicates whether the check passed
          const passed = Object.values(result).find(val => typeof val === 'boolean' && val !== result.isLoading);

          if (passed) {
            return true;
          } else {
            // Keep the returnUrl in the query params
            return this.router.createUrlTree([redirectTo], {
              queryParams: { returnUrl: state.url }
            });
          }
        })
      );
    };
  }
}

/**
 * Guard to check if user has an active trial subscription
 */
export const canActivateTrial: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  return inject(SubscriptionGuard).createGuard(selectHasActiveTrial, '/upgrade')(route, state);
};

/**
 * Guard to check if user has a basic subscription
 */
export const canActivateBasic: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  return inject(SubscriptionGuard).createGuard(selectHasActiveBasic, '/upgrade')(route, state);
};

/**
 * Guard to check if user has a pro subscription
 */
export const canActivatePro: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  return inject(SubscriptionGuard).createGuard(selectHasActivePro, '/upgrade')(route, state);
};

/**
 * Guard to check if user has at least a basic subscription (basic or pro)
 */
export const canActivatePaid: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  return inject(SubscriptionGuard).createGuard(selectHasActivePaid, '/upgrade')(route, state);
};

/**
 * Guard to check if user has any active subscription
 */
export const canActivateAnySubscription: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  return inject(SubscriptionGuard).createGuard(selectIsSubscriptionActive, '/renew-subscription')(route, state);
};
