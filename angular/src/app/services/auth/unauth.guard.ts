import { inject, Injectable } from '@angular/core';
import {
  ActivatedRouteSnapshot,
  RouterStateSnapshot,
  Router,
  CanActivateFn,
} from '@angular/router';

import { select, Store } from '@ngrx/store';
import { Observable, tap } from 'rxjs';
import { AppState } from 'src/app/store/interfaces/app-state';
import { isAnonymousSelector } from 'src/app/store/selectors/auth.selectors';

/** UnAuthentication Guard, only unauthenticated users with no JWT Token on localstorage */
@Injectable({
  providedIn: 'root',
})
export class UnAuthGuard {
  constructor(private store: Store<AppState>, public router: Router) { }

  canActivate(
    route: ActivatedRouteSnapshot,
    state: RouterStateSnapshot
  ): Observable<boolean> {
    return this.store.pipe(
      select(isAnonymousSelector),
      tap((isAnon) => {
        if (!isAnon) {
          this.router.navigateByUrl('/app');
        }
      })
    );
  }
}

export const canActivateAnon: CanActivateFn = (route: ActivatedRouteSnapshot, state: RouterStateSnapshot) => {
  return inject(UnAuthGuard).canActivate(route, state);
};
