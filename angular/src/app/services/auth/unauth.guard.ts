import { Injectable } from '@angular/core';
import {
  CanActivate,
  ActivatedRouteSnapshot,
  RouterStateSnapshot,
  Router,
} from '@angular/router';

import { select, Store } from '@ngrx/store';
import { map, Observable, tap } from 'rxjs';
import { AppState } from 'src/app/store/interfaces/app-state';
import { isAnonymousSelector } from 'src/app/store/selectors/auth.selectors';

@Injectable({
  providedIn: 'root',
})
/** UnAuthentication Guard, only unauthenticated users with no JWT Token on localstorage */
@Injectable()
export class UnAuthGuard implements CanActivate {
  constructor(private store: Store<AppState>, public router: Router) {}

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
