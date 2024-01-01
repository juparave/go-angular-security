import { HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { catchError, map, of, switchMap, tap } from 'rxjs';
import { AuthService } from 'src/app/services/auth/auth.service';
import { PersistanceService } from 'src/app/services/persistance.service';
import { User } from 'src/app/model/user';
import {
  registerAction,
  registerSuccessAction,
  registerFailureAction,
} from 'src/app/store/actions/auth.actions';

@Injectable()
export class RegisterEffect {
  constructor(
    private actions$: Actions,
    private authService: AuthService,
    private persistanceService: PersistanceService,
    private router: Router
  ) { }

  register$ = createEffect(() =>
    this.actions$.pipe(
      ofType(registerAction),
      switchMap((action) => {
        return this.authService.register(action.request).pipe(
          map((currentUser: User) => {
            this.persistanceService.set('accessToken', currentUser.accessToken);
            this.persistanceService.set(
              'refreshToken',
              currentUser.refreshToken
            );
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
    )
  );

  redirecAfterSubmit$ = createEffect(
    () =>
      this.actions$.pipe(
        ofType(registerSuccessAction),
        tap(() => {
          this.router.navigateByUrl('/app');
        })
      ),
    // doesn't return an Observable, so we set dispatch to false
    { dispatch: false }
  );
}
