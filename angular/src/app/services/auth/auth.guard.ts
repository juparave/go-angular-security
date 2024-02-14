import { inject, Injectable } from '@angular/core';
import {
  ActivatedRouteSnapshot,
  RouterStateSnapshot,
  Router,
  CanActivateFn,
} from '@angular/router';

import { select, Store } from '@ngrx/store';
import { filter, map, Observable, tap } from 'rxjs';
import { AppState } from 'src/app/store/interfaces/app-state';
import { selectIsLoggedIn } from 'src/app/store/selectors/auth.selectors';

/** Authentication Guard, only authenticated users with JWT Token on localstorage */
@Injectable({
  providedIn: 'root',
})
export class AuthGuard {
  constructor(private store: Store<AppState>, public router: Router) { }

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): Observable<boolean> {
    return this.store.select(selectIsLoggedIn).pipe(
      filter(s => !s.isLoading),
      tap((s) => {
        if (!s.isLoggedIn) {
          this.router.navigate(['/login'], { queryParams: { returnUrl: state.url } })
        }
      }),
      map((s) => s.isLoggedIn ?? false)
    );
  }
}

export const canActivateAuth: CanActivateFn = (route: ActivatedRouteSnapshot, state: RouterStateSnapshot) => {
  return inject(AuthGuard).canActivate(route, state);
};
