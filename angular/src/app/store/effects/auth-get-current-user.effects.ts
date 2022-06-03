import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { catchError, map, of, switchMap, tap } from 'rxjs';
import { AuthService } from 'src/app/services/auth/auth.service';

import { PersistanceService } from 'src/app/services/persistance.service';
import { User } from 'src/app/model/user';

import {
  getCurrentUserAction,
  getCurrentUserSuccessAction,
  getCurrentUserFailureAction,
} from 'src/app/store/actions/auth.actions';

@Injectable()
export class GetCurrentUserEffect {
  constructor(
    private actions$: Actions,
    private authService: AuthService,
    private persistanceService: PersistanceService,
    private router: Router
  ) {}

  getCurrentUser$ = createEffect(() =>
    this.actions$.pipe(
      ofType(getCurrentUserAction),
      switchMap(() => {
        const token = this.persistanceService.get('accessToken');

        if (!token) {
          // if no token on localStorage
          return of(getCurrentUserFailureAction());
        }

        return this.authService.getCurrentUser().pipe(
          map((currentUser: User) => {
            // this.persistanceService.set('accessToken', currentUser.token);
            return getCurrentUserSuccessAction({ currentUser });
          }),
          catchError(() => {
            return of(getCurrentUserFailureAction());
          })
        );
      })
    )
  );

  getCurrentUserFailure$ = createEffect(
    () =>
      this.actions$.pipe(
        ofType(getCurrentUserFailureAction),
        tap(() => {
          // this.router.navigateByUrl('/login');
        })
      ),
    // doesn't return an Observable, so we set dispatch to false
    { dispatch: false }
  );
}
