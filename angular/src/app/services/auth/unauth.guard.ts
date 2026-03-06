import { inject, Injectable } from '@angular/core';
import {
  ActivatedRouteSnapshot,
  RouterStateSnapshot,
  Router,
  CanActivateFn,
} from '@angular/router';

import { Store } from '@ngrx/store';
import { Observable, tap, map, filter } from 'rxjs';
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
      filter(state => !state.isLoading),
      tap((state) => {
        if (!state.isAnon) {
          this.router.navigateByUrl('/app');
        }
      }),
      map((state) => state.isAnon),
    );
  }
}

export const canActivateAnon: CanActivateFn = (route: ActivatedRouteSnapshot, state: RouterStateSnapshot) => {
  return inject(UnAuthGuard).canActivate(route, state);
};
