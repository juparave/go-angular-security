import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

/// ref: https://blog.angular-university.io/angular-loading-indicator/

@Injectable({
  providedIn: 'root'
})
export class LoadingService {
  private loadingSubject = new BehaviorSubject<boolean>(false);

  loading$ = this.loadingSubject.asObservable();

  loadingOn() {
    this.loadingSubject.next(true);
  }

  loadingOff() {
    this.loadingSubject.next(false);
  }
}
