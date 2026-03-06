import { HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { catchError, map, of, switchMap, tap } from 'rxjs';
import { AuthService } from 'src/app/services/auth/auth.service';
import { PersistanceService } from 'src/app/services/persistance.service';
import { User } from 'src/app/models/user';
import {
  loginAction,
  loginFailureAction,
  loginSuccessAction,
} from 'src/app/store/actions/auth.actions';

@Injectable()
export class LoginEffect {
  constructor(
    private actions$: Actions,
    private authService: AuthService,
    private persistanceService: PersistanceService,
    private router: Router
  ) { }

  login$ = createEffect(() => {
    return this.actions$.pipe(
      ofType(loginAction),
      switchMap((action) => {
        return this.authService.login(action.request).pipe(
          map((currentUser: User) => {
            this.persistanceService.set('token', currentUser.accessToken);
            this.persistanceService.set(
              'refreshToken',
              currentUser.refreshToken
            );
            return loginSuccessAction({ currentUser, returnUrl: action.returnUrl });
          }),
          catchError((errorResponse: HttpErrorResponse) => {
            // console.log(errorResponse);
            return of(
              loginFailureAction({ errors: errorResponse.error.errors })
            );
          })
        );
      })
    )
  });

  redirecAfterSubmit$ = createEffect(
    () => {
      return this.actions$.pipe(
        ofType(loginSuccessAction),
        tap((action) => {
          this.router.navigateByUrl(action.returnUrl);
        })
      );
    },
    // doesn't return an Observable, so we set dispatch to false
    { dispatch: false }
  );
}
