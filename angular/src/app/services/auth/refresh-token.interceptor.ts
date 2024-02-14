/// ref: https://www.bezkoder.com/angular-12-refresh-token/
import { Injectable } from '@angular/core';
import {
  HttpInterceptor,
  HttpEvent,
  HttpHandler,
  HttpRequest,
  HttpErrorResponse,
} from '@angular/common/http';
import { BehaviorSubject, Observable, throwError } from 'rxjs';
import { catchError, filter, switchMap, take, tap } from 'rxjs/operators';
import { Store } from '@ngrx/store';
import { selectIsLoggedIn } from 'src/app/store/selectors/auth.selectors';
import { AuthService } from './auth.service';
import {
  refreshTokenAction,
  refreshTokenFailureAction,
  refreshTokenSuccessAction,
} from 'src/app/store/actions/auth.actions';
import { PersistanceService } from 'src/app/services/persistance.service';
import { AppState } from 'src/app/store/interfaces/app-state';
import { User } from 'src/app/model/user';

const TOKEN_HEADER_KEY = 'Token';

@Injectable()
export class RefreshTokenInterceptor implements HttpInterceptor {
  private isRefreshing = false;
  private refreshTokenSubject: BehaviorSubject<any> = new BehaviorSubject<any>(
    null
  );

  constructor(
    private store: Store<AppState>,
    private persistanceService: PersistanceService,
    private authService: AuthService
  ) { }

  intercept(
    req: HttpRequest<any>,
    next: HttpHandler
  ): Observable<HttpEvent<Object>> {
    let authReq = req;
    const token = this.persistanceService.get('accessToken');
    if (token != null) {
      authReq = this.addTokenHeader(req, token);
    }

    return next.handle(authReq).pipe(
      catchError((error) => {
        if (
          error instanceof HttpErrorResponse &&
          !authReq.url.includes('login') &&
          error.status === 401
        ) {
          return this.handle401Error(authReq, next);
        }

        return throwError(() => error);
      })
    );
  }

  private handle401Error(request: HttpRequest<any>, next: HttpHandler) {
    if (!this.isRefreshing) {
      this.isRefreshing = true;
      this.refreshTokenSubject.next(null);
      this.store.dispatch(refreshTokenAction());

      const refreshToken = this.persistanceService.get('refreshToken');

      if (refreshToken)
        return this.authService.refreshTokens(refreshToken).pipe(
          switchMap((user: User) => {
            console.log('Refresh token success');
            this.isRefreshing = false;

            this.store.dispatch(
              refreshTokenSuccessAction({ currentUser: user })
            );

            this.refreshTokenSubject.next(user.accessToken);

            return next.handle(this.addTokenHeader(request, user.accessToken));
          }),
          catchError((err) => {
            console.log('Refresh token failed');
            this.isRefreshing = false;

            // this.authService.logout().then(() => console.log("Logout"));
            this.store.dispatch(
              refreshTokenFailureAction({
                errors: { ['error']: ['Refresh token failed'] },
              })
            );
            return throwError(() => err);
          })
        );
    }

    return this.refreshTokenSubject.pipe(
      tap((token) => {
        console.log('Refresh token: ', token);
      }),
      filter((token) => token !== null),
      take(1),
      switchMap((token) => next.handle(this.addTokenHeader(request, token)))
    );
  }

  private addTokenHeader(request: HttpRequest<any>, token: string) {
    /* for Spring Boot back-end */
    // return request.clone({ headers: request.headers.set(TOKEN_HEADER_KEY, 'Bearer ' + token) });

    /* for Node.js Express back-end */
    return request.clone({
      headers: request.headers.set(TOKEN_HEADER_KEY, token),
    });
  }
}
