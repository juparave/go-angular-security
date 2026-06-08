import {
  HttpContextToken,
  HttpEvent,
  HttpHandler,
  HttpInterceptor,
  HttpRequest,
} from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { LoadingService } from './loading.service';
import { Observable, finalize } from 'rxjs';

/// ref: https://blog.angular-university.io/angular-loading-indicator/

export const SkipLoading = new HttpContextToken<boolean>(() => false);

@Injectable()
export class LoadingInterceptor implements HttpInterceptor {
  private loadingService = inject(LoadingService);

  intercept(req: HttpRequest<unknown>, next: HttpHandler): Observable<HttpEvent<unknown>> {
    // Check for a custom attribute
    // to avoid showing loading spinner
    if (req.context.get(SkipLoading)) {
      // Pass the request directly to the next handler
      return next.handle(req);
    }

    // Turn on the loading spinner
    this.loadingService.loadingOn();

    return next.handle(req).pipe(
      finalize(() => {
        // Turn off the loading spinner
        this.loadingService.loadingOff();
      })
    );
  }
}
