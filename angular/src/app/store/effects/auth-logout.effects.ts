import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { Actions, createEffect, ofType } from '@ngrx/effects';
import { map, tap } from 'rxjs';
import { PersistanceService } from 'src/app/services/persistance.service';
import {
  logoutAction,
  logoutSuccessAction,
} from 'src/app/store/actions/auth.actions';

@Injectable()
export class LogoutEffect {
  constructor(
    private actions$: Actions,
    private persistanceService: PersistanceService,
    private router: Router
  ) { }

  logout$ = createEffect(() => {
    return this.actions$.pipe(
      ofType(logoutAction),
      map(() => {
        // clear token from persistanceService
        this.persistanceService.remove('token');
        this.persistanceService.remove('refreshToken');

        return logoutSuccessAction();
      }),
    )
  });

  redirecAfterLogoutSuccess$ = createEffect(
    () => {
      return this.actions$.pipe(
        ofType(logoutSuccessAction),
        tap(() => {
          this.router.navigateByUrl('/login');
        })
      );
    },
    // doesn't return an Observable, so we set dispatch to false
    { dispatch: false }
  );
}
