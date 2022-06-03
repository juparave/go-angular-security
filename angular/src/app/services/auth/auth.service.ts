// services/auth/auth.service.ts
import { Injectable } from '@angular/core';
import { Router } from '@angular/router';
import {
  HttpClient,
  HttpErrorResponse,
  HttpHeaders,
  HttpResponse,
} from '@angular/common/http';

import {
  catchError,
  map,
  Observable,
  of,
  shareReplay,
  take,
  tap,
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
  private serviceUrl = environment.apiPath + '/';
  headers = new HttpHeaders().set(
    'Content-Type',
    'application/json; charset=utf-8'
  );
  constructor(private http: HttpClient, public router: Router) {}

  /**
   * Calls Register service
   * @param data: User
   */
  register(data: User): Observable<User> {
    const url = environment.apiPath + '/register';
    return this.http.post<User>(url, data);
  }

  /**
   * Calls login service, if succesful, loads and stores sessions values
   * @param data: LoginRequest
   */
  login(data: LoginRequest): Observable<User> {
    const url = environment.apiPath + '/login';
    return this.http.post<any>(url, data.user).pipe(
      map((res) => res.user),
      shareReplay()
    );
  }

  getCurrentUser(): Observable<User> {
    const url = environment.apiPath + '/user';
    return this.http.get<User>(url);
  }

  /**
   * RefreshTokens from server, user must be logged in
   */
  public refreshTokens(refreshToken: string) {
    const url = environment.apiPath + '/refresh_token';
    let headers = new HttpHeaders();
    headers = headers.set('Content-Type', 'application/json; charset=utf-8');
    console.log('auth service get tokens...');

    return this.http
      .post<any>(
        url,
        {
          refreshToken: refreshToken,
        },
        { headers: headers }
      )
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

  private handleErrors(errorResponse: HttpErrorResponse) {
    console.error('Error: ' + JSON.stringify(errorResponse));
    return throwError(() => errorResponse.error.errors);
  }
}
