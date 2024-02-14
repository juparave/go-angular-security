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
import { selectIsAnonymous } from 'src/app/store/selectors/auth.selectors';

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
    return this.store.select(selectIsAnonymous).pipe(
      filter(s => !s.isLoading),
      tap((s) => {
        if (!s.isAnon) {
          this.router.navigateByUrl('/app');
        }
      }),
      map((s) => s.isAnon ?? false)
    );
  }
}

export const canActivateAnon: CanActivateFn = (route: ActivatedRouteSnapshot, state: RouterStateSnapshot) => {
  return inject(UnAuthGuard).canActivate(route, state);
};
