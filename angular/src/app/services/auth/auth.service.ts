// services/auth/auth.service.ts
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import { HttpClient, HttpErrorResponse, HttpHeaders } from '@angular/common/http';

import { map, Observable, shareReplay, throwError } from 'rxjs';

import { environment } from 'src/environments/environment';
import { LoginRequest } from 'src/app/store/types/login-request.interface';
import { User } from 'src/app/models/user';
import { RequestResetPassword, ResetPassword } from 'src/app/store/types/request-reset.interface';
import { AccountRegistration } from 'src/app/store/types/account-registration.interface';

export interface ChangePasswordRequest {
  currentPassword: string;
  newPassword: string;
  confirmPassword: string;
}

@Injectable({
  providedIn: 'root',
})
export class AuthService {
  private baseUrl = `${environment.apiUrl}`;
  private headers = new HttpHeaders().set('Content-Type', 'application/json; charset=utf-8');

  constructor(
    private http: HttpClient,
    public router: Router
  ) {}

  /**
   * Register a new user (legacy method)
   * @param data: User
   */
  register(data: User): Observable<User> {
    return this.http.post<User>(`${this.baseUrl}/register`, data);
  }

  /**
   * Register a new account with first user as admin
   * @param data: RegisterAccountRequest
   */
  registerAccount(data: AccountRegistration): Observable<{ message: string; user: User }> {
    return this.http.post<{ message: string; user: User }>(
      `${this.baseUrl}/register-account`,
      data
    );
  }

  /**
   * Login service
   * @param data: LoginRequest
   */
  login(data: LoginRequest): Observable<User> {
    return this.http.post<{ user: User }>(`${this.baseUrl}/login`, data.user).pipe(
      map((res) => res.user),
      shareReplay()
    );
  }

  /**
   * Gets the current logged-in user
   */
  getCurrentUser(): Observable<User> {
    return this.http.get<User>(`${this.baseUrl}/user`).pipe(
      map((user) => {
        if (user === null) {
          throw new HttpErrorResponse({ status: 401, statusText: 'unauthenticated' });
        }
        return user;
      })
    );
  }

  /**
   * Refreshes tokens from server
   * @param refreshToken The refresh token
   */
  refreshTokens(refreshToken: string): Observable<User> {
    return this.http
      .post<{
        user: User;
      }>(`${this.baseUrl}/refresh-token`, { refreshToken }, { headers: this.headers })
      .pipe(map((res) => res.user));
  }

  /**
   * Logout method
   */
  logout(): Observable<{ message: string }> {
    return this.http.post<{ message: string }>(`${this.baseUrl}/logout`, {});
  }

  /**
   * Request password reset
   * @param data Request reset data with email
   */
  requestReset(data: RequestResetPassword): Observable<{ message: string }> {
    return this.http.post<{ message: string }>(`${this.baseUrl}/request-password-reset`, data);
  }

  /**
   * Reset password with token
   * @param data Reset data with token and new password
   */
  resetPassword(data: ResetPassword & { token: string }): Observable<{ message: string }> {
    return this.http.post<{ message: string }>(`${this.baseUrl}/reset-password`, data);
  }

  /**
   * Change password for authenticated user
   * @param data Change password data
   */
  changePassword(data: ChangePasswordRequest): Observable<{ message: string }> {
    return this.http.post<{ message: string }>(`${this.baseUrl}/change-password`, data);
  }

  /**
   * Checks if a JWT token is expired
   * @param token The JWT token to check
   * @returns true if token is expired, false otherwise
   */
  isTokenExpired(token: string): boolean {
    try {
      const expiry = JSON.parse(atob(token.split('.')[1])).exp;
      return Math.floor(new Date().getTime() / 1000) >= expiry;
    } catch {
      return true;
    }
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
