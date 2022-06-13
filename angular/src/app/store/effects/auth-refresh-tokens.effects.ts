import { HttpErrorResponse } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { catchError, map, of, switchMap, tap } from 'rxjs';
import { AuthService } from 'src/app/services/auth/auth.service';
import { PersistanceService } from 'src/app/services/persistance.service';
import { User } from 'src/app/model/user';
import {
  refreshTokenAction,
  refreshTokenSuccessAction,
  refreshTokenFailureAction,
} from 'src/app/store/actions/auth.actions';

@Injectable()
export class RefreshTokensEffect {
  constructor(
    private actions$: Actions,
    private authService: AuthService,
    private persistanceService: PersistanceService,
    private router: Router
  ) {}

  persistNewTokens$ = createEffect(
    () =>
      this.actions$.pipe(
        ofType(refreshTokenSuccessAction),
        tap((action) => {
          this.persistanceService.set(
            'accessToken',
            action.currentUser.accessToken
          );
          this.persistanceService.set(
            'refreshToken',
            action.currentUser.refreshToken
          );
        })
      ),
    { dispatch: false }
  );

  redirecAfterFailure$ = createEffect(
    () =>
      this.actions$.pipe(
        ofType(refreshTokenFailureAction),
        tap(() => {
          // clear token from persistanceService
          this.persistanceService.remove('accessToken');
          this.persistanceService.remove('refreshToken');

          this.router.navigateByUrl('/login');
        })
      ),
    // doesn't return an Observable, so we set dispatch to false
    { dispatch: false }
  );
}
