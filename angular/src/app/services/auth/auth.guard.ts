import { inject, Injectable } from '@angular/core';
import {
  ActivatedRouteSnapshot,
  RouterStateSnapshot,
  Router,
  CanActivateFn,
} from '@angular/router';

import { select, Store } from '@ngrx/store';
import { map, Observable, tap } from 'rxjs';
import { AppState } from 'src/app/store/interfaces/app-state';
import { isLoggedInSelector } from 'src/app/store/selectors/auth.selectors';

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
    return this.store.pipe(
      select(isLoggedInSelector),
      tap((loggedIn) => {
        if (!loggedIn) {
          this.router.navigateByUrl('/login');
        }
      }),
      map((loggedIn) => loggedIn ?? false)
    );
  }
}

export const canActivateAuth: CanActivateFn = (route: ActivatedRouteSnapshot, state: RouterStateSnapshot) => {
  return inject(AuthGuard).canActivate(route, state);
};
