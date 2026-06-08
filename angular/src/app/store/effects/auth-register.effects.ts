import { HttpErrorResponse } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Router } from '@angular/router';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { catchError, map, of, switchMap, tap } from 'rxjs';
import { AuthService } from 'src/app/services/auth/auth.service';
import { PersistanceService } from 'src/app/services/persistance.service';
import { User } from 'src/app/models/user';
import {
  registerAction,
  registerSuccessAction,
  registerFailureAction,
  registerAccountAction,
  registerAccountSuccessAction,
  registerAccountFailureAction,
} from 'src/app/store/actions/auth.actions';

@Injectable()
export class RegisterEffect {
  private actions$ = inject(Actions);
  private authService = inject(AuthService);
  private persistanceService = inject(PersistanceService);
  private router = inject(Router);

  register$ = createEffect(() => {
    return this.actions$.pipe(
      ofType(registerAction),
      switchMap((action) => {
        return this.authService.register(action.request).pipe(
          map((currentUser: User) => {
            this.persistanceService.set('token', currentUser.accessToken);
            this.persistanceService.set('refreshToken', currentUser.refreshToken);
            return registerSuccessAction({ currentUser });
          }),
          catchError((errorResponse: HttpErrorResponse) => {
            console.log('register error response: ', errorResponse);
            return of(
              registerFailureAction({
                errors: { register: [errorResponse.error.message] },
              })
            );
          })
        );
      })
    );
  });

  registerAccount$ = createEffect(() => {
    return this.actions$.pipe(
      ofType(registerAccountAction),
      switchMap((action) => {
        return this.authService.registerAccount(action.request).pipe(
          map((response) => {
            return registerAccountSuccessAction({ currentUser: response.user });
          }),
          catchError((errorResponse: HttpErrorResponse) => {
            return of(
              registerAccountFailureAction({
                errors: { register: [errorResponse.error?.message ?? 'Registration failed'] },
              })
            );
          })
        );
      })
    );
  });

  redirecAfterSubmit$ = createEffect(
    () => {
      return this.actions$.pipe(
        ofType(registerSuccessAction, registerAccountSuccessAction),
        tap(() => {
          this.router.navigateByUrl('/app');
        })
      );
    },
    { dispatch: false }
  );
}
