import { inject, Injectable } from '@angular/core';
import { CanActivateFn, Router, UrlTree, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { Observable, map, filter } from 'rxjs';
import { Store } from '@ngrx/store';
import { AppState } from 'src/app/store/interfaces/app-state';
import { selectIsSubscriptionActive } from 'src/app/store/selectors/subscription.selectors';

/**
 * Guard to prevent users with an active subscription from accessing certain routes (e.g., /subscription)
 */
@Injectable({
  providedIn: 'root',
})
export class NoActiveSubscriptionGuard {
  constructor(private store: Store<AppState>, public router: Router) { }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean | UrlTree> {
    return this.store.select(selectIsSubscriptionActive).pipe(
      filter(result => !result.isLoading),
      map(result => {
        console.log("SubscriptionGuard", result);
        if (result.isActive) {
          // Redirect to /app if subscription is active
          return this.router.createUrlTree(['/app']);
        } else {
          return true;
        }
      })
    );
  }
}

export const canActivateNoActiveSubscription: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
) => {
  return inject(NoActiveSubscriptionGuard).canActivate(route, state);
};
