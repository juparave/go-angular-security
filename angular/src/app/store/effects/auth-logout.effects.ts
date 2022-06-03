import { HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { catchError, map, of, switchMap, tap } from 'rxjs';
import { AuthService } from 'src/app/services/auth/auth.service';
import { PersistanceService } from 'src/app/services/persistance.service';
import { User } from 'src/app/model/user';
import {
  logoutAction,
  logoutFailureAction,
  logoutSuccessAction,
} from 'src/app/store/actions/auth.actions';

@Injectable()
export class LogoutEffect {
  constructor(
    private actions$: Actions,
    private authService: AuthService,
    private persistanceService: PersistanceService,
    private router: Router
  ) {}

  logout$ = createEffect(
    () =>
      this.actions$.pipe(
        ofType(logoutAction),
        tap(() => {
          // action to do on logout
        })
      ),
    { dispatch: false }
  );

  redirecAfterLogoutSuccess$ = createEffect(
    () =>
      this.actions$.pipe(
        ofType(logoutSuccessAction),
        tap(() => {
          this.router.navigateByUrl('/app');
        })
      ),
    // doesn't return an Observable, so we set dispatch to false
    { dispatch: false }
  );
}
