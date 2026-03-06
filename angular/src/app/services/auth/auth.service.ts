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
import { User } from 'src/app/models/user';
import {
  RequestResetPassword,
  ResetPassword,
  ValidateRequestResetPassword,
} from 'src/app/store/types/request-reset.interface';

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private baseUrl = `${environment.apiUrl}`;
  private headers = new HttpHeaders().set(
    'Content-Type',
    'application/json; charset=utf-8'
  );

  constructor(private http: HttpClient, public router: Router) { }

  /**
   * Calls Register service
   * @param data: User
   */
  register(data: User): Observable<User> {
    return this.http.post<User>(`${this.baseUrl}/register`, data);
  }

  /**
   * Calls login service, if succesful, loads and stores sessions values
   * @param data: LoginRequest
   */
  login(data: LoginRequest): Observable<User> {
    return this.http.post<any>(`${this.baseUrl}/login`, data.user).pipe(
      map((res) => res.user),
      shareReplay()
    );
  }

  /**
   * Gets the current logged-in user
   */
  getCurrentUser(): Observable<User> {
    return this.http.get<User>(`${this.baseUrl}/user`).pipe(
      map(user => {
        if (user === null) {
          throw new HttpErrorResponse({ status: 401, statusText: 'unauthenticated' });
        }
        return user;
      })
    );
  }

  /**
   * Refreshes tokens from server, user must be logged in
   * @param refreshToken The refresh token
   */
  refreshTokens(refreshToken: string): Observable<User> {
    return this.http.post<any>(
      `${this.baseUrl}/refresh_token`,
      { refreshToken },
      { headers: this.headers }
    ).pipe(map((res) => res.user));
  }

  /**
   * Logout method
   * Note: Implementation delegated to NgRx store
   */
  logout(): Observable<never> {
    return throwError(() => new Error('Not implemented, use store for logout'));
  }

  /**
   * Request password reset
   * @param data Request reset data with email
   */
  requestReset(data: RequestResetPassword): Observable<any> {
    return this.http.post<any>(
      `${this.baseUrl}/request_reset_password`,
      data
    );
  }

  /**
   * Validate reset code for password reset
   * @param data Validation data with email and code
   */
  validateRequestReset(data: ValidateRequestResetPassword): Observable<any> {
    return this.http.post<any>(
      `${this.baseUrl}/validate_request_reset_password`,
      data
    );
  }

  /**
   * Reset password with new credentials
   * @param data Reset data with email and new password
   */
  resetPassword(data: ResetPassword): Observable<any> {
    return this.http.post<any>(`${this.baseUrl}/reset_password`, data);
  }

  /**
   * Checks if a JWT token is expired
   * @param token The JWT token to check
   * @returns true if token is expired, false otherwise
   */
  isTokenExpired(token: string): boolean {
    const expiry = (JSON.parse(atob(token.split('.')[1]))).exp;
    return (Math.floor((new Date).getTime() / 1000)) >= expiry;
  }

  /**
   * Handles HTTP errors
   * @param errorResponse The HTTP error response
   * @returns Observable with error
   */
  private handleErrors(errorResponse: HttpErrorResponse) {
    console.error('Error: ' + JSON.stringify(errorResponse));
    return throwError(() => errorResponse.error.errors);
  }
}