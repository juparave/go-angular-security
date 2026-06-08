import { Injectable, inject } from '@angular/core';
import { Router } from '@angular/router';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { catchError, of, switchMap, tap } from 'rxjs';
import { AuthService } from 'src/app/services/auth/auth.service';

import { PersistanceService } from 'src/app/services/persistance.service';
import { User } from 'src/app/models/user';

import {
  getCurrentUserAction,
  getCurrentUserSuccessAction,
  getCurrentUserFailureAction,
  refreshTokenAction,
} from 'src/app/store/actions/auth.actions';
import { getSubscriptionAction } from 'src/app/store/actions/subscription.actions';
import { Store } from '@ngrx/store';

@Injectable()
export class GetCurrentUserEffect {
  private store = inject(Store);
  private actions$ = inject(Actions);
  private authService = inject(AuthService);
  private persistanceService = inject(PersistanceService);
  private router = inject(Router);

  getCurrentUser$ = createEffect(() => {
    return this.actions$.pipe(
      ofType(getCurrentUserAction),
      switchMap(() => {
        const token = this.persistanceService.get('token');

        if (!token) {
          // if no token on localStorage
          return of(getCurrentUserFailureAction());
        }

        if (this.authService.isTokenExpired(token)) {
          // if token expired, try to refresh it by triggering refreshTokenAction
          // this.store.dispatch(refreshTokenAction());
        }

        return this.authService.getCurrentUser().pipe(
          // Use switchMap to dispatch multiple actions
          switchMap((currentUser: User) => {
            // Dispatch both success action and the action to get subscription
            return of(
              getCurrentUserSuccessAction({ currentUser }),
              getSubscriptionAction() // Dispatch this after user success
            );
          }),
          catchError(() => {
            return of(getCurrentUserFailureAction());
          })
        );
      })
    );
  });

  getCurrentUserFailure$ = createEffect(
    () => {
      return this.actions$.pipe(
        ofType(getCurrentUserFailureAction),
        tap(() => {
          // clear token from persistanceService
          this.persistanceService.remove('token');
          this.persistanceService.remove('refreshToken');
          // redirection is handled by the router guard
          // this.router.navigateByUrl('/login');
        })
      );
    },
    // doesn't return an Observable, so we set dispatch to false
    { dispatch: false }
  );

  doRefreshToken$ = createEffect(
    () => {
      return this.actions$.pipe(
        ofType(refreshTokenAction),
        tap(() => {
          // this is pending
        })
      );
    },
    { dispatch: false }
  );
}
