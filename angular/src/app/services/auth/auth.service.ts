// services/auth/auth.service.ts
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import {
  HttpClient,
  HttpErrorResponse,
  HttpHeaders,
} from '@angular/common/http';

import {
  map,
  Observable,
  shareReplay,
  throwError,
} from 'rxjs';

import { environment } from 'src/environments/environment';
import { LoginRequest } from 'src/app/store/types/login-request.interface';
import { User } from 'src/app/model/user';
import {
  RequestResetPassword,
  ResetPassword,
  ValidateRequestResetPassword,
} from 'src/app/store/types/request-reset.interface';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private serviceUrl = environment.apiUrl + '/';
  headers = new HttpHeaders().set(
    'Content-Type',
    'application/json; charset=utf-8'
  );
  constructor(private http: HttpClient, public router: Router) { }

  /**
   * Calls Register service
   * @param data: User
   */
  register(data: User): Observable<User> {
    const url = environment.apiUrl + '/register';
    return this.http.post<User>(url, data);
  }

  /**
   * Calls login service, if succesful, loads and stores sessions values
   * @param data: LoginRequest
   */
  login(data: LoginRequest): Observable<User> {
    const url = environment.apiUrl + '/login';
    return this.http.post<any>(url, data.user).pipe(
      map((res) => res.user),
      shareReplay()
    );
  }

  getCurrentUser(): Observable<User> {
    const url = environment.apiUrl + '/user';
    // check if user is not null
    return this.http.get<User>(url).pipe(
      map(user => {
        if (user === null) {
          throw new HttpErrorResponse({ status: 401, statusText: 'unauthenticated' });
        }
        return user;
      })
    );
  }

  /**
   * RefreshTokens from server, user must be logged in
   */
  public refreshTokens(refreshToken: string) {
    const url = environment.apiUrl + '/refresh_token';
    let headers = new HttpHeaders();
    headers = headers.set('Content-Type', 'application/json; charset=utf-8');
    console.log('auth service get tokens...');

    return this.http.post<any>(url, { refreshToken: refreshToken, }, { headers: headers })
      .pipe(map((res) => res.user));
  }

  public logout() {
    throwError(() => Error('Not implemented, use store for logout'));
  }

  public requestReset(data: RequestResetPassword) {
    return this.http.post<any>(
      `${this.serviceUrl}request_reset_password`,
      data
    );
  }

  public validateRequestReset(data: ValidateRequestResetPassword) {
    return this.http.post<any>(
      `${this.serviceUrl}validate_request_reset_password`,
      data
    );
  }

  public resetPassword(data: ResetPassword) {
    return this.http.post<any>(`${this.serviceUrl}reset_password`, data);
  }

  public isTokenExpired(token: string): boolean {
    const expiry = (JSON.parse(atob(token.split('.')[1]))).exp;
    return (Math.floor((new Date).getTime() / 1000)) >= expiry;
  }

  private handleErrors(errorResponse: HttpErrorResponse) {
    console.error('Error: ' + JSON.stringify(errorResponse));
    return throwError(() => errorResponse.error.errors);
  }
}
